package models

import (
	"fmt"
	"sync"
	"time"
)

type IssueType int

const (
	Unknown IssueType = iota
	Payment
	MutualFund
	Gold
	Insurance
)

func (it IssueType) String() string {
	switch it {
	case Payment:
		return "Payment"
	case MutualFund:
		return "Mutual Fund"
	case Gold:
		return "Gold"
	case Insurance:
		return "Insurance"
	default:
		return "Unknown"
	}
}

type IssueStatus int

const (
	Created IssueStatus = iota
	InProgress
	Resolved
)

func (it IssueStatus) String() string {
	switch it {
	case Created:
		return "Created"
	case InProgress:
		return "InProgress"
	case Resolved:
		return "Resolved"
	default:
		return "Unknown"
	}
}

type Issue struct {
	Id          string      `json:"id"`
	TxnId       string      `json:"txn_id"`
	Type        IssueType   `json:"type"`
	Subject     string      `json:"subject"`
	Description string      `json:"description"`
	Email       string      `json:"email"`
	Status      IssueStatus `json:"status"`
	Resolution  string      `json:"resolution"`
	mu          sync.RWMutex
	// additional fields to track metadata of the issue
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewIssue(id, txnId, subject, description, email string, issueType IssueType) (*Issue, error) {
	if txnId == "" || subject == "" || description == "" || email == "" {
		return nil, fmt.Errorf("invalid input")
	}
	return &Issue{
		Id:          id,
		TxnId:       txnId,
		Type:        issueType,
		Subject:     subject,
		Description: description,
		Email:       email,
		Status:      Created,
		CreatedAt:   time.Now().Unix(),
	}, nil
}

// to get the status of the issue
func (i *Issue) GetStatus() IssueStatus {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Status
}

func (i *Issue) UpdateStatus(status IssueStatus, resolution string) (bool, error) {
	if resolution == "" {
		return false, fmt.Errorf("resolution cannot be empty")
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Status = status
	i.UpdatedAt = time.Now().Unix()
	i.Resolution = resolution
	return true, nil
}
