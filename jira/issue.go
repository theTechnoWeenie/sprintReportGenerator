package jira

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Value struct {
	Points float64 `json:"value"`
}

type Estimate struct {
	FieldId string `json:"statFieldId"`
	Value   Value  `json:"statFieldValue"`
}

func (e *Estimate) GetPoints() float64 {
	return e.Value.Points
}

type WorkItem struct {
	SecondsLogged int `json:"timeSpentSeconds"`
}

type Worklog struct {
	Log []WorkItem `json:"worklogs"`
}

func (w *Worklog) GetSecondsWorked() int {
	sum := 0
	for _, item := range w.Log {
		sum += item.SecondsLogged
	}
	return sum
}

type Issue struct {
	Id          int      `josn:"id"`
	Key         string   `json:"key"`
	Type        string   `json:"typeName"`
	StoryPoints Estimate `json:"estimateStatistic"`
	Status      string   `json:"statusName"`
	hoursLogged int
}

func (i *Issue) GetEstimate() float64 {
	return i.StoryPoints.GetPoints()
}

func (i *Issue) GetIssueMarkup() string {
	return fmt.Sprintf("{jiraissues:key=%s}", i.Key)
}

func (i *Issue) GetWorkLogged(client Client) int {
	if i.hoursLogged != 0 {
		return i.hoursLogged
	}
	worklogJson := client.Get(fmt.Sprintf("rest/api/2/issue/%s/worklog", i.Key), nil)
	var worklog Worklog
	err := json.Unmarshal(worklogJson, &worklog)
	if err != nil {
		log.Fatal("Could not parse worklog json")
	}
	i.hoursLogged = worklog.GetSecondsWorked()
	return i.hoursLogged
}

type IssueList []Issue

func (i *IssueList) GetTotalEstimate(addedIssueKeys []string) float64 {
	sum := 0.0
	for _, issue := range *i {
		//TODO: this should be reworked to use the filter function.
		shouldFilter := false
		for _, added := range addedIssueKeys {
			if shouldFilter || added == issue.Key {
				shouldFilter = true
			}
		}
		if !shouldFilter {
			sum += issue.GetEstimate()
		}
	}
	return sum
}

func (i *IssueList) GetTotalHours(closed bool, client Client) string {
	sum := 0
	for _, issue := range *i {
		if !closed || strings.ToLower(issue.Status) == "closed" || strings.ToLower(issue.Status) == "resolved" {
			sum += issue.GetWorkLogged(client)
		}
	}
	return fmt.Sprintf("%d h", sum/60/60)
}

func (list *IssueList) FilterIssues(filterKeys []string, returnFiltered bool) IssueList {
	filtered := IssueList{}
	for _, issue := range *list {
		shouldFilter := false
		for _, key := range filterKeys {
			shouldFilter = shouldFilter || (issue.Key == key)
		}
		if shouldFilter == returnFiltered {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func (list *IssueList) Append(otherList IssueList) IssueList {
	ret := *list
	for _, issue := range otherList {
		ret = append(ret, issue)
	}
	return ret
}

func (list *IssueList) GetIssueKeys() []string {
	keys := []string{}
	for _, issue := range *list {
		keys = append(keys, issue.Key)
	}
	return keys
}
