package service

import (
	"fmt"
	"iss/internal/models"
	"os"
	"sync"
)

type ResolutionService struct {
	issueService  *IssueService
	AgentService  *AgentService
	strategy      AssignmentStrategy
	issueAgentMap map[string]string
	mutex         sync.RWMutex
}

func NewResolutionService(issueService *IssueService, agentService *AgentService, strategy AssignmentStrategy) *ResolutionService {
	if strategy == nil {
		strategy = &FreeAgentFirstStrategy{}
	}
	return &ResolutionService{
		issueService:  issueService,
		AgentService:  agentService,
		strategy:      strategy,
		issueAgentMap: make(map[string]string),
	}
}

func (rs *ResolutionService) CreateIssue(txnID, subject, description, email string, issueType models.IssueType) (string, error) {
	return rs.issueService.CreateIssue(txnID, subject, description, email, issueType)
}

func (rs *ResolutionService) AddAgent(email, name string, expertise map[models.IssueType]bool) (string, error) {
	return rs.AgentService.AddAgent(email, name, expertise)
}

func (rs *ResolutionService) AssignIssue(issueId string) (string, bool, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	waitListed := false

	issue := rs.issueService.GetIssue(issueId)
	if issue == nil {
		return "", waitListed, fmt.Errorf("issue not found")
	}

	targetAgent := rs.strategy.Assign(issue, rs.AgentService.GetAvailableAgentsByExpertise(), rs.AgentService.GetBusyAgentHeap())
	if targetAgent == nil {
		fmt.Println(rs.AgentService.GetAvailableAgentsByExpertise())
		fmt.Println(rs.AgentService.GetBusyAgentHeap())
		os.Exit(0)
	}
	waitListed, err := rs.AgentService.AssignIssue(targetAgent, issue)
	if err != nil {
		return "", waitListed, fmt.Errorf("error occurred - assign issue %w", err)
	}
	if !waitListed {
		rs.issueAgentMap[issueId] = targetAgent.Id
	}
	return targetAgent.Id, waitListed, nil
}

func (rs *ResolutionService) GetIssues(filter map[string]string) []*models.Issue {
	return rs.issueService.GetIssues(filter)
}

func (rs *ResolutionService) UpdateIssue(issueId, resolution string, status models.IssueStatus) error {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	if _, ok := rs.issueAgentMap[issueId]; !ok {
		return fmt.Errorf("cannot update, issue not yet assigned to any agent")
	}
	return rs.issueService.UpdateIssue(issueId, resolution, status)
}

func (rs *ResolutionService) ResolveIssue(issueId, resolution string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	issue := rs.issueService.GetIssue(issueId)
	if issue == nil {
		return fmt.Errorf("issue not found")
	}

	agentId := rs.issueAgentMap[issueId]
	newIssueAssigned, err := rs.AgentService.ResolveIssue(agentId, resolution)
	if err != nil {
		fmt.Println("error occurred - ResolveIssue", err)
		return err
	}
	if newIssueAssigned != nil {
		rs.issueAgentMap[newIssueAssigned.Id] = agentId
		fmt.Printf("Issue %s has been assigned to agent %s from the pendingIssues Queue \n", newIssueAssigned.Id, agentId)
	}

	err = rs.issueService.UpdateIssue(issueId, resolution, models.Resolved)
	if err != nil {
		fmt.Println("error occurred - UpdateIssue", err)
		return err
	}

	return nil
}

func (rs *ResolutionService) ViewAgentsWorkHistory() map[string][]string {
	return rs.AgentService.GetWorkHistory()
}
