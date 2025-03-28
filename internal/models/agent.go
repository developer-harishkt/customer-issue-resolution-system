package models

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Agent struct {
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Expertise      map[IssueType]bool `json:"expertise"`
	AssignedIssue  *Issue
	PendingIssues  []*Issue          // considering it as a list to assume the issues would be picked up in FIFO Order
	ResolvedIssues map[string]*Issue // stores the resolved issues by their ID
	HeapIndex      int
	CreatedAt      int64 `json:"created_at"`
	mu             sync.RWMutex
}

func NewAgent(id, name, email string, expertise map[IssueType]bool) (*Agent, error) {
	if name == "" || email == "" || len(expertise) == 0 {
		return nil, fmt.Errorf("invalid agent data")
	}

	return &Agent{
		Id:             id,
		Name:           name,
		Email:          email,
		Expertise:      expertise,
		PendingIssues:  []*Issue{},
		ResolvedIssues: make(map[string]*Issue),
		CreatedAt:      time.Now().Unix(),
	}, nil
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

func (a *Agent) GetExpertise() map[IssueType]bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Expertise
}

func (a *Agent) GetPendingIssues() []*Issue {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.PendingIssues
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

func (a *Agent) AddToPendingIssues(issue *Issue) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.PendingIssues = append(a.PendingIssues, issue)
	fmt.Printf("Issue %s has been added to pending-issues-queue of agent %s \n", issue.Id, a.Id)
}

func (a *Agent) AssignIssue(issue *Issue) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.AssignedIssue != nil {
		return errors.New("agent is already assigned an issue")
	}
	a.AssignedIssue = issue
	return nil
}

// resolve issue automatically assigns an issue if there are any pending issues
func (a *Agent) ResolveIssue(resolution string) (*Issue, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.AssignedIssue != nil {
		a.ResolvedIssues[a.AssignedIssue.Id] = a.AssignedIssue
		a.AssignedIssue = nil
		if len(a.PendingIssues) > 0 {
			a.AssignedIssue = a.PendingIssues[0]
			a.PendingIssues = a.PendingIssues[1:]
		}
		return a.AssignedIssue, nil
	} else {
		return nil, fmt.Errorf("no assigned issue to resolve")
	}
}
