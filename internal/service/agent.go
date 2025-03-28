package service

import (
	"fmt"
	m "iss/internal/models"
	"sync"
	"sync/atomic"
)

type AgentService struct {
	agents          map[string]*m.Agent
	availableAgents int32
	idCounter       int32
	mu              sync.RWMutex
}

func NewAgentService() *AgentService {
	return &AgentService{
		agents:          make(map[string]*m.Agent),
		availableAgents: 0,
	}
}

func (as *AgentService) AddAgent(email, name string, expertise []m.IssueType) string {
	as.mu.Lock()
	defer as.mu.Unlock()

	id := fmt.Sprintf("A%d", atomic.AddInt32(&as.idCounter, 1))
	expertiseMap := make(map[m.IssueType]bool)
	for _, it := range expertise {
		expertiseMap[it] = true
	}
	agent := m.NewAgent(id, email, name, expertiseMap)
	as.agents[id] = agent
	atomic.AddInt32(&as.availableAgents, 1)
	return id
}

func (as *AgentService) GetAgent(id string) *m.Agent {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return as.agents[id]
}

func (as *AgentService) GetAgents() []*m.Agent {
	as.mu.RLock()
	defer as.mu.RUnlock()

	agents := make([]*m.Agent, 0, len(as.agents))
	for _, agent := range as.agents {
		agents = append(agents, agent)
	}
	return agents
}

func (as *AgentService) AssignIssue(agentID string, issue *m.Issue, toPending bool) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if agent, exists := as.agents[agentID]; exists {
		if toPending {
			agent.AddToPendingIssues(issue)
		} else {
			if agent.IsAvailable() {
				agent.AssignIssue(issue)
				as.availableAgents--
			}
		}
	}
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

func (s *AgentService) GetAvailableAgents() int32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.availableAgents
}

func (as *AgentService) ResolveIssue(agentID, resolution string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if agent, exists := as.agents[agentID]; exists {
		if !agent.IsAvailable() {
			agent.ResolveIssue(resolution)
			if agent.GetAssignedIssue() == nil && len(agent.GetPendingIssues()) == 0 {
				atomic.AddInt32(&as.availableAgents, 1)
			}
		}
	}
}
