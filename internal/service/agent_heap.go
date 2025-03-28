package service

import (
	"container/heap"
	"iss/internal/models"
)

type AgentHeap []*models.Agent

func (h AgentHeap) Len() int           { return len(h) }
func (h AgentHeap) Less(i, j int) bool { return len(h[i].PendingIssues) < len(h[j].PendingIssues) }
func (h AgentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].HeapIndex, h[j].HeapIndex = i, j
}

func (h *AgentHeap) Push(x interface{}) {
	agent := x.(*models.Agent)
	agent.HeapIndex = len(*h)
	*h = append(*h, agent)
}

func (h *AgentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	agent := old[n-1]
	agent.HeapIndex = -1
	*h = old[0 : n-1]
	return agent
}

func InitializeHeap() *AgentHeap {
	h := &AgentHeap{}
	heap.Init(h)
	return h
}
