package service

import (
	"container/heap"
	"iss/internal/models"
)

type AssignmentStrategy interface {
	Assign(issue *models.Issue, availableAgentsMapByExpertise map[models.IssueType]map[string]*models.Agent, busyAgentHeap *AgentHeap) *models.Agent
}

type AssignmentStrategies int

const (
	FreeAgentFirst AssignmentStrategies = iota
)

type FreeAgentFirstStrategy struct{}

func (s *FreeAgentFirstStrategy) Assign(issue *models.Issue, availableAgentsMapByExpertise map[models.IssueType]map[string]*models.Agent, busyAgentHeap *AgentHeap) *models.Agent {
	var minQueueAgent *models.Agent
	expertise := issue.Type

	// if an agent with desired expertise is available
	if len(availableAgentsMapByExpertise[expertise]) > 0 {
		for _, agent := range availableAgentsMapByExpertise[expertise] {
			return agent
		}
	}

	// check if any agent with different expertise
	for _, agentMap := range availableAgentsMapByExpertise {
		if len(agentMap) > 0 {
			for _, agent := range agentMap {
				return agent
			}
		}
	}

	// if all the agents are busy
	if busyAgentHeap.Len() > 0 {
		minQueueAgent = heap.Pop(busyAgentHeap).(*models.Agent)
	}

	return minQueueAgent
}

func NewFreeAgentFirstStrategy() *FreeAgentFirstStrategy {
	return &FreeAgentFirstStrategy{}
}

func GetAssignmentStrategy(as AssignmentStrategies) AssignmentStrategy {
	switch as {
	default:
		return NewFreeAgentFirstStrategy()
	}
}
