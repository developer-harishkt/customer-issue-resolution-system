package service

import (
	"fmt"
	m "iss/internal/models"
	"strings"
	"sync"
)

type IssueService struct {
	Issues map[string]*m.Issue
	mu     sync.RWMutex
}

func NewIssueService() *IssueService {
	return &IssueService{
		Issues: make(map[string]*m.Issue),
	}
}

func (is *IssueService) CreateIssue(txnID, subject, description, email string, issueType m.IssueType) string {
	is.mu.Lock()
	defer is.mu.Unlock()

	id := "I" + txnID // in ideal systems we should be using uuid's
	issue := m.NewIssue(id, txnID, subject, description, email, issueType)
	is.Issues[id] = issue
	return id
}

func (is *IssueService) GetIssue(id string) *m.Issue {
	is.mu.RLock()
	defer is.mu.RUnlock()
	return is.Issues[id]
}

func (is *IssueService) UpdateIssue(issueId, resolution string, status m.IssueStatus) (bool, error) {
	is.mu.Lock()
	defer is.mu.Unlock()
	if issue, exists := is.Issues[issueId]; exists {
		issue.UpdateStatus(status, resolution)
		return true, nil
	} else {
		return false, fmt.Errorf("issue with id %s not found", issueId)
	}
}

func (is *IssueService) GetIssues(filter map[string]string) []*m.Issue {
	is.mu.RLock()
	defer is.mu.RUnlock()

	var filteredIssues []*m.Issue
	for _, issue := range is.Issues {
		found := true
		for key, value := range filter {
			switch strings.ToLower(key) {
			case "id":
				if issue.Id != value {
					found = false
				}
			case "txnid":
				if issue.TxnId != value {
					found = false
				}
			case "type":
				if issue.Type.String() != value {
					found = false
				}
			case "subject":
				if issue.Subject != value {
					found = false
				}
			case "description":
				if issue.Description != value {
					found = false
				}
			case "email":
				if issue.Email != value {
					found = false
				}
			case "status":
				if string(issue.Status) != value {
					found = false
				}
			case "resolution":
				if issue.Resolution != value {
					found = false
				}
			default:
				continue
			}
			if !found {
				break
			}
		}
		if found {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	return filteredIssues
}
