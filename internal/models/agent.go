package models

import (
	"errors"
	"sync"
)

type Agent struct {
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Expertise      map[IssueType]bool `json:"expertise"`
	AssignedIssue  *Issue
	PendingIssues  []*Issue          // considering it as a list to assume the issues would be picked up in FIFO Order
	ResolvedIssues map[string]*Issue // stores the resolved issues by their ID
	CreatedAt      int64             `json:"created_at"`
	mu             sync.RWMutex
}

func NewAgent(id, name, email string, expertise map[IssueType]bool) *Agent {
	return &Agent{
		Id:             id,
		Name:           name,
		Email:          email,
		Expertise:      expertise,
		PendingIssues:  []*Issue{},
		ResolvedIssues: make(map[string]*Issue),
	}
}

func (a *Agent) IsAvailable() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.AssignedIssue == nil
}

func (a *Agent) GetAssignedIssue() *Issue {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.AssignedIssue
}

func (a *Agent) AssignIssue(issue *Issue) (bool, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.AssignedIssue != nil {
		// if the agent is already assigned an issue, return false and an error
		return false, errors.New("agent is already assigned an issue")
	}
	a.AssignedIssue = issue
	return true, nil
}

func (a *Agent) GetPendingIssues() []*Issue {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.PendingIssues
}

func (a *Agent) AddToPendingIssues(issue *Issue) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.PendingIssues = append(a.PendingIssues, issue)
}

// resolve issue automatically assigns an issue if there are any pending issues
func (a *Agent) ResolveIssue(resolution string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.AssignedIssue != nil {
		a.ResolvedIssues[a.AssignedIssue.Id] = a.AssignedIssue
		a.AssignedIssue = nil
		if len(a.PendingIssues) > 0 {
			a.AssignedIssue = a.PendingIssues[0]
			a.PendingIssues = a.PendingIssues[1:]
		}
	}
}

func (a *Agent) GetResolvedIssues() map[string]*Issue {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.ResolvedIssues
}

func (a *Agent) HasExpertise(it IssueType) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Expertise[it]
}
