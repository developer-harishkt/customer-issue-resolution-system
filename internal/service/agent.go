package service

import (
	"container/heap"
	"fmt"
	m "iss/internal/models"
	"sync"
	"sync/atomic"
)

type AgentService struct {
	agents                     map[string]*m.Agent
	AvailableAgentsByExpertise map[m.IssueType]map[string]*m.Agent
	busyAgentHeap              *AgentHeap
	idCounter                  int32
	mu                         sync.RWMutex
}

func NewAgentService() *AgentService {
	return &AgentService{
		agents:                     make(map[string]*m.Agent),
		AvailableAgentsByExpertise: make(map[m.IssueType]map[string]*m.Agent),
		busyAgentHeap:              InitializeHeap(),
	}
}

func (as *AgentService) AddAgent(email, name string, expertise map[m.IssueType]bool) (string, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	id := fmt.Sprintf("A%d", atomic.AddInt32(&as.idCounter, 1))
	agent, err := m.NewAgent(id, email, name, expertise)
	if err != nil {
		fmt.Println("error occurred", err)
		return "", err
	}
	as.agents[id] = agent
	for expertise, _ := range agent.Expertise {
		if _, ok := as.AvailableAgentsByExpertise[expertise]; ok {
			as.AvailableAgentsByExpertise[expertise][agent.Id] = agent
		} else {
			as.AvailableAgentsByExpertise[expertise] = map[string]*m.Agent{
				agent.Id: agent,
			}
		}
	}
	return id, nil
}

func (as *AgentService) GetAgent(id string) *m.Agent {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.agents[id]
}

func (as *AgentService) GetAgents() map[string]*m.Agent {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.agents
}

func (as *AgentService) GetAvailableAgentsByExpertise() map[m.IssueType]map[string]*m.Agent {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.AvailableAgentsByExpertise
}

func (as *AgentService) GetBusyAgentHeap() *AgentHeap {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.busyAgentHeap
}

func (as *AgentService) AssignIssue(agent *m.Agent, issue *m.Issue) (bool, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	waitListed := true
	if agent.IsAvailable() {
		err := agent.AssignIssue(issue)
		if err != nil {
			return false, fmt.Errorf("error occurred %w", err)
		}
		for expertise, _ := range agent.GetExpertise() {
			delete(as.AvailableAgentsByExpertise[expertise], agent.Id)
		}
		waitListed = false
		heap.Push(as.busyAgentHeap, agent)
	} else {
		agent.AddToPendingIssues(issue)
		if agent.HeapIndex >= 0 && agent.HeapIndex < len(*as.busyAgentHeap) && (*as.busyAgentHeap)[agent.HeapIndex] == agent {
			heap.Push(as.busyAgentHeap, agent)
		}
	}
	return waitListed, nil
}

func (as *AgentService) GetWorkHistory() map[string][]string {
	as.mu.RLock()
	defer as.mu.RUnlock()

	history := make(map[string][]string)
	for _, agent := range as.agents {
		resolved := agent.GetResolvedIssues()
		issueIDs := make([]string, 0, len(resolved))
		for id := range resolved {
			issueIDs = append(issueIDs, id)
		}
		history[agent.Id] = issueIDs
	}
	return history
}

func (as *AgentService) ResolveIssue(agentId, resolution string) (*m.Issue, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if agent, ok := as.agents[agentId]; ok {
		newIssueAssigned, err := agent.ResolveIssue(resolution)
		if err != nil {
			return nil, err
		}
		if newIssueAssigned != nil {
			heap.Push(as.busyAgentHeap, agent)
		} else {
			for expertise, _ := range agent.GetExpertise() {
				as.AvailableAgentsByExpertise[expertise] = map[string]*m.Agent{
					agentId: agent,
				}
			}
		}
		return newIssueAssigned, nil
	} else {
		return nil, fmt.Errorf("agent not found")
	}
}
