package main

import (
	"fmt"
	"iss/internal/models"
	"iss/internal/service"
)

func main() {
	issueService := service.NewIssueService()
	agentService := service.NewAgentService()
	resolutionService := service.NewResolutionService(issueService, agentService, nil)

	fmt.Println("Creating issues...")
	i1 := resolutionService.CreateIssue("T1", "Payment Failed", "My payment failed but money is debited", "testUser1@test.com", models.Payment)
	fmt.Printf("Issue %s created against transaction T1\n", i1)

	i2 := resolutionService.CreateIssue("T2", "Purchase Failed", "Unable to purchase Mutual Fund", "testUser2@test.com", models.MutualFund)
	fmt.Printf("Issue %s created against transaction T2\n", i2)

	i3 := resolutionService.CreateIssue("T3", "Payment Failed", "My payment failed but money is debited", "testUser2@test.com", models.Payment)
	fmt.Printf("Issue %s created against transaction T3\n", i3)

	fmt.Println("\nAdding agents...")
	a1 := resolutionService.AddAgent("agent1@test.com", "Agent 1", []models.IssueType{models.Payment, models.Gold})
	fmt.Printf("Agent %s created\n", a1)

	a2 := resolutionService.AddAgent("agent2@test.com", "Agent 2", []models.IssueType{models.Payment})
	fmt.Printf("Agent %s created\n", a2)

	fmt.Println("\nAssigning issues...")
	resolutionService.AssignIssue(i1)
	fmt.Printf("Issue %s assigned to agent %s\n", i1, a1)

	resolutionService.AssignIssue(i2)
	fmt.Printf("Issue %s assigned to agent %s\n", i2, a2)

	resolutionService.AssignIssue(i3)
	fmt.Printf("Issue %s added to waitlist of Agent %s\n", i3, a1)

	fmt.Println("\nGetting issues for testUser2@test.com...")
	issues := resolutionService.GetIssues(map[string]string{"email": "testUser2@test.com"})
	for _, issue := range issues {
		fmt.Printf("%v\n", issue)
	}

	fmt.Println("\nUpdating and resolving issues...")
	resolutionService.UpdateIssue(i1, "Waiting for payment confirmation", models.InProgress)
	fmt.Printf("Issue %s status updated to In Progress\n", i1)

	err := resolutionService.ResolveIssue(i1, "Payment reversed")
	if err != nil {
		fmt.Println("error occured while resolving the issue", err)
	} else {
		fmt.Printf("Issue %s resolved\n", i1)
	}

	err = resolutionService.ResolveIssue(i2, "PaymentFailed debited amount will get reversed")
	if err != nil {
		fmt.Println("error occured while resolving the issue", err)
	} else {
		fmt.Printf("Issue %s resolved\n", i2)
	}

	fmt.Println("\nViewing agents' work history...")
	history := resolutionService.ViewAgentsWorkHistory()
	for agentID, resolvedIssues := range history {
		fmt.Printf("%s -> %v\n", agentID, resolvedIssues)
	}
}
