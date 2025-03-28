package main

import (
	"fmt"
	"iss/internal/models"
	"iss/internal/service"
	"os"
)

func main() {
	issueService := service.NewIssueService()
	agentService := service.NewAgentService()
	assignmentStrategy := service.GetAssignmentStrategy(service.FreeAgentFirst)
	resolutionService := service.NewResolutionService(issueService, agentService, assignmentStrategy)

	// scenario 1: 8 tasks for 4 agents, 4 picked, 4 queued
	fmt.Println("scenario 1: 8 tasks for 4 agents, 4 picked, 4 queued")
	testIssues := []models.Issue{
		{TxnId: "T1", Subject: "Payment Failed", Description: "My payment failed but money is debited", Email: "testUser1@test.com", Type: models.Payment},
		{TxnId: "T2", Subject: "Purchase Failed", Description: "Unable to purchase Mutual Fund", Email: "testUser2@test.com", Type: models.MutualFund},
		{TxnId: "T3", Subject: "Payment Failed", Description: "My payment failed but money is debited", Email: "testUser2@test.com", Type: models.Payment},
		{TxnId: "T4", Subject: "Purchase Failed", Description: "Unable to purchase Mutual Fund", Email: "testUser1@test.com", Type: models.Payment},
		{TxnId: "T5", Subject: "Payment Failed", Description: "My payment failed but money is debited", Email: "testUser1@test.com", Type: models.Payment},
		{TxnId: "T6", Subject: "Purchase Failed", Description: "Unable to purchase Mutual Fund", Email: "testUser2@test.com", Type: models.MutualFund},
		{TxnId: "T7", Subject: "Payment Failed", Description: "My payment failed but money is debited", Email: "testUser1@test.com", Type: models.Payment},
		{TxnId: "T8", Subject: "Purchase Failed", Description: "Unable to purchase Mutual Fund", Email: "testUser2@test.com", Type: models.MutualFund},
	}
	activeIssueIds := make([]string, 0)
	fmt.Println("Creating issues...")
	for i := range testIssues {
		issue := &testIssues[i]
		id, err := resolutionService.CreateIssue(issue.TxnId, issue.Subject, issue.Description, issue.Email, issue.Type)
		if err != nil {
			fmt.Println("Error occurred - CreateIssue:", err)
			os.Exit(1)
		}
		fmt.Printf("Issue %s created against transaction %s\n", id, issue.TxnId)
		activeIssueIds = append(activeIssueIds, id)
	}

	testAgents := []models.Agent{
		{Email: "agent1@test.com", Name: "Agent 1", Expertise: map[models.IssueType]bool{models.Payment: true}},
		{Email: "agent2@test.com", Name: "Agent 2", Expertise: map[models.IssueType]bool{models.Payment: true}},
		{Email: "agent3@test.com", Name: "Agent 3", Expertise: map[models.IssueType]bool{models.MutualFund: true}},
		{Email: "agent4@test.com", Name: "Agent 4", Expertise: map[models.IssueType]bool{models.MutualFund: true}},
	}
	fmt.Println("\nAdding agents...")
	for i := range testAgents {
		agent := &testAgents[i]
		id, err := resolutionService.AddAgent(agent.Email, agent.Name, agent.Expertise)
		if err != nil {
			fmt.Println("Error Occurred - AddAgent:", err)
		}
		fmt.Printf("Agent %s created\n", id)
	}

	fmt.Println(resolutionService.AgentService.AvailableAgentsByExpertise)

	fmt.Println("assigning issues...")
	for _, issueId := range activeIssueIds {
		agentId, waitlisted, err := resolutionService.AssignIssue(issueId)
		if err != nil {
			fmt.Println("error occurred - AssignIssue", issueId)
		} else if waitlisted {
			fmt.Printf("Issue %s added to waitlist of Agent %s\n", issueId, agentId)
		} else {
			fmt.Printf("Issue %s assigned to Agent %s\n", issueId, agentId)
		}
		if issueId == "IT2" {
			fmt.Println(resolutionService.AgentService.AvailableAgentsByExpertise)
		}
	}

	fmt.Println()
	fmt.Println()
	// scenario 2: A1 should have min queue length & top of heap, hence T9 should be added to pending list of A1
	fmt.Println("scenario 2: A1 should have min queue length & top of heap, hence T9 should be added to pending list of A1")
	issueId := "IT1"
	err := resolutionService.ResolveIssue(issueId, "random-message")
	if err != nil {
		fmt.Printf("Error resolving issue %s: %v\n", issueId, err)
	} else {
		fmt.Printf("Issue %s resolved\n", issueId)
	}

	id, err := resolutionService.CreateIssue("T9", "a", "a", "a", models.Payment)
	if err != nil {
		fmt.Println("Error occurred - CreateIssue:", err)
		os.Exit(1)
	}

	agentId, waitlisted, err := resolutionService.AssignIssue(id)
	if err != nil {
		fmt.Println("error occurred - AssignIssue", id)
	} else if waitlisted {
		fmt.Printf("Issue %s added to waitlist of Agent %s\n", id, agentId)
	} else {
		fmt.Printf("Issue %s assigned to Agent %s\n", id, agentId)
	}

	fmt.Println()
	fmt.Println()
	// scenario 3: after resolving all the issues for A2, A2 should be in avaialbleAgentByExpertise
	fmt.Println("scenario 3: after resolving all the issues for A2, A2 should be in avaialbleAgentByExpertise")
	issueId = "IT3"
	err = resolutionService.ResolveIssue(issueId, "random-message")
	if err != nil {
		fmt.Printf("Error resolving issue %s: %v\n", issueId, err)
	} else {
		fmt.Printf("Issue %s resolved\n", issueId)
	}

	issueId = "IT7"
	err = resolutionService.ResolveIssue(issueId, "random-message")
	if err != nil {
		fmt.Printf("Error resolving issue %s: %v\n", issueId, err)
	} else {
		fmt.Printf("Issue %s resolved\n", issueId)
	}

	id, err = resolutionService.CreateIssue("T10", "a", "a", "a", models.Payment)
	if err != nil {
		fmt.Println("Error occurred - CreateIssue:", err)
		os.Exit(1)
	}

	agentId, waitlisted, err = resolutionService.AssignIssue(id)
	if err != nil {
		fmt.Println("error occurred - AssignIssue", id)
	} else if waitlisted {
		fmt.Printf("Issue %s added to waitlist of Agent %s\n", id, agentId)
	} else {
		fmt.Printf("Issue %s assigned to Agent %s\n", id, agentId)
	}

	fmt.Println("\nGetting issues for testUser2@test.com")
	issues := resolutionService.GetIssues(map[string]string{"email": "testUser2@test.com"})
	for _, issue := range issues {
		fmt.Printf("%v\n", issue)
	}

	fmt.Println("\nViewing agents' work history...")
	history := resolutionService.ViewAgentsWorkHistory()
	for agentID, resolvedIssues := range history {
		fmt.Printf("%s -> %v\n", agentID, resolvedIssues)
	}
}
