package service

import (
	"fmt"
	"iss/internal/models"
	"sync"
)

type AssignmentStrategy interface {
	Assign(issue *models.Issue, agents []*models.Agent) *models.Agent
}

type FreeAgentFirstStrategy struct{}

func (s *FreeAgentFirstStrategy) Assign(issue *models.Issue, agents []*models.Agent) *models.Agent {
	var freeAgent *models.Agent
	var minQueueAgent *models.Agent
	minQueueSize := -1

	for _, agent := range agents {
		if agent.HasExpertise(issue.Type) {
			if agent.GetAssignedIssue() == nil {
				freeAgent = agent
				break
			} else if minQueueSize == -1 || len(agent.GetPendingIssues()) < minQueueSize {
				minQueueSize = len(agent.GetPendingIssues())
				minQueueAgent = agent
			}
		}
	}

	if freeAgent != nil {
		return freeAgent
	}

	return minQueueAgent
}

type ResolutionService struct {
	issueService *IssueService
	agentService *AgentService
	strategy     AssignmentStrategy
	mutex        sync.RWMutex
}

func NewResolutionService(issueService *IssueService, agentService *AgentService, strategy AssignmentStrategy) *ResolutionService {
	if strategy == nil {
		strategy = &FreeAgentFirstStrategy{}
	}
	return &ResolutionService{
		issueService: issueService,
		agentService: agentService,
		strategy:     strategy,
	}
}

func (rs *ResolutionService) CreateIssue(txnID, subject, description, email string, issueType models.IssueType) string {
	return rs.issueService.CreateIssue(txnID, subject, description, email, issueType)
}

func (rs *ResolutionService) AddAgent(email, name string, expertise []models.IssueType) string {
	return rs.agentService.AddAgent(email, name, expertise)
}

func (rs *ResolutionService) AssignIssue(issueID string) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	issue := rs.issueService.GetIssue(issueID)
	if issue == nil {
		return
	}

	agents := rs.agentService.GetAgents()
	targetAgent := rs.strategy.Assign(issue, agents)

	if targetAgent != nil {
		toPending := targetAgent.GetAssignedIssue() != nil
		rs.agentService.AssignIssue(targetAgent.Id, issue, toPending)
	}

}

func (rs *ResolutionService) GetIssues(filter map[string]string) []*models.Issue {
	return rs.issueService.GetIssues(filter)
}

func (rs *ResolutionService) UpdateIssue(issueID, resolution string, status models.IssueStatus) {
	rs.issueService.UpdateIssue(issueID, resolution, status)
}

func (rs *ResolutionService) ResolveIssue(issueID, resolution string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	issue := rs.issueService.GetIssue(issueID)
	if issue == nil {
		return fmt.Errorf("issue not found")
	}

	for _, agent := range rs.agentService.GetAgents() {
		if assigned := agent.GetAssignedIssue(); assigned != nil && assigned.Id == issueID {
			rs.agentService.ResolveIssue(agent.Id, resolution)
			rs.issueService.UpdateIssue(issueID, resolution, models.Resolved)
			return nil
		}
	}
	return fmt.Errorf("issue not currently assigned to any agent")
}

func (rs *ResolutionService) ViewAgentsWorkHistory() map[string][]string {
	return rs.agentService.GetWorkHistory()
}
