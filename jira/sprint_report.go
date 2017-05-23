package jira

import (
	"fmt"
	"strings"
)

type EstimateTotals struct {
	Total float64 `json:"value"`
}

type SprintContents struct {
	CompleteIssues       IssueList       `json:"completedIssues"`
	IncompleteIssues     IssueList       `json:"issuesNotCompletedInCurrentSprint"`
	RemovedIssues        IssueList       `json:"puntedIssues"`
	TotalPointsCompelted EstimateTotals  `json:"completedIssuesEstimateSum"`
	TotalPointsMissed    EstimateTotals  `json:"issuesNotCompletedEstimateSum"`
	TotalPointsRemoved   EstimateTotals  `json:"puntedIssuesEstimateSum"`
	TotalCommited        EstimateTotals  `json:"allIssuesEstimateSum"`
	AddedIssues          map[string]bool `json:"issueKeysAddedDuringSprint"`
}

type Report struct {
	Contents      SprintContents `json:"contents"`
	SprintDetails Sprint         `json:"sprint"`
}

func (r *Report) GetAddedIssueKeys() []string {
	addedIssues := make([]string, 0, len(r.Contents.AddedIssues))
	for k := range r.Contents.AddedIssues {
		addedIssues = append(addedIssues, k)
	}
	return addedIssues
}

func (r *Report) GetEmergentIssues(onlyCompleted bool) IssueList {
	emergentIssues := IssueList{}
	for _, issue := range r.Contents.CompleteIssues {
		if strings.ToLower(issue.Type) == "emergent task" || strings.ToLower(issue.Type) == "bug" {
			emergentIssues = append(emergentIssues, issue)
		}
	}
	if !onlyCompleted {
		for _, issue := range r.Contents.IncompleteIssues {
			if strings.ToLower(issue.Type) == "emergent task" || strings.ToLower(issue.Type) == "bug" {
				emergentIssues = append(emergentIssues, issue)
			}
		}
	}
	return emergentIssues
}

func (r *Report) GetCommitedIssues() IssueList {
	added := r.GetAddedIssueKeys()
	emergent := r.GetEmergentIssues(true)
	for _, issue := range emergent {
		added = append(added, issue.Key)
	}
	commited := r.Contents.CompleteIssues.FilterIssues(added, false)
	commited = commited.Append(r.Contents.IncompleteIssues.FilterIssues(added, false))
	return commited
}

func (r *Report) GenerateCommitedTable() string {
	summary := "h2. Sprint Commitment: \n\n"
	summary += "|| User Story || Estimate ||\n"
	for _, i := range r.GetCommitedIssues() {
		summary += fmt.Sprintf("| %s | %d |\n", i.GetIssueMarkup(), int(i.GetEstimate()))
	}
	summary += "\n\n"
	return summary
}

func (r *Report) GenerateAdditionsTable() string {
	addedIssueKeys := r.GetAddedIssueKeys()
	addedIssues := r.Contents.CompleteIssues.FilterIssues(addedIssueKeys, true)
	for _, issue := range r.Contents.IncompleteIssues.FilterIssues(addedIssueKeys, true) {
		addedIssues = append(addedIssues, issue)
	}
	emergentIssues := r.GetEmergentIssues(false)
	addedIssues = addedIssues.FilterIssues(emergentIssues.GetIssueKeys(), false)
	summary := "h2. Sprint Additions: \n\n"
	summary += "|| User Story || Estimate ||\n"
	for _, issue := range addedIssues {
		summary += fmt.Sprintf("| %s | %d |\n", issue.GetIssueMarkup(), int(issue.GetEstimate()))
	}
	summary += "\n\n"
	return summary
}

func (r *Report) GenerateEmergentTable(client Client) string {
	summary := "h2. Emergent Work\n\n"
	summary += "|| Work Item: || Type ||\n"
	for _, issue := range r.GetEmergentIssues(true) {
		summary += fmt.Sprintf("| %s | %s |\n", issue.GetIssueMarkup(), issue.Type)
	}
	summary += "\n\n"
	return summary
}

func (r *Report) GenerateSummaryTable(client Client, velocity Velocity) string {
	summary := "h2. Sprint Details\n\n"
	summary += fmt.Sprintf("| Start Date: | %s |\n", r.SprintDetails.StartDate)
	summary += fmt.Sprintf("| End Date: | %s |\n", r.SprintDetails.EndDate)
	commitedPoints := r.Contents.CompleteIssues.GetTotalEstimate(r.GetAddedIssueKeys()) + r.Contents.IncompleteIssues.GetTotalEstimate(r.GetAddedIssueKeys())
	totalPoints := r.Contents.CompleteIssues.GetTotalEstimate(nil) + r.Contents.IncompleteIssues.GetTotalEstimate(nil)
	summary += fmt.Sprintf("| Committed Points:  | %d |\n", int(commitedPoints))
	summary += fmt.Sprintf("| Added Points: | %d |\n", int(totalPoints-commitedPoints))
	summary += fmt.Sprintf("| Removed Points: | %d |\n", int(r.Contents.RemovedIssues.GetTotalEstimate(nil)))
	summary += fmt.Sprintf("| Completed Points: | %d |\n", int(r.Contents.CompleteIssues.GetTotalEstimate(nil)))
	numToAverage := 3
	summary += fmt.Sprintf("| Average Velocity (last %d): | %d |\n", numToAverage, int(velocity.GetAverageVelocity(numToAverage)))
	summary += "\n\n"
	return summary
}

func (r *Report) GenerateReport(client Client, velocity Velocity) string {
	summary := r.GenerateSummaryTable(client, velocity)
	summary += r.GenerateCommitedTable()
	summary += r.GenerateEmergentTable(client)
	summary += r.GenerateAdditionsTable()
	return summary
}
