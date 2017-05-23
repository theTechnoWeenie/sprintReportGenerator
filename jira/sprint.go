package jira

import "fmt"

type Sprint struct {
	Id               int    `json:"id"`
	Sequence         int    `json:"sequence"`
	Name             string `json:"name"`
	State            string `json:"state"`
	LinkedPagesCount int    `json:"linkedPagesCount"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	CompleteDate     string `json:"completeDate"`
}

type SprintList struct {
	Sprints   []Sprint `json:"sprints"`
	ProjectId int      `json:"rapidViewId"`
}

func (s *SprintList) FindSprintWithName(name string) (Sprint, error) {
	for _, sprint := range s.Sprints {
		if name == sprint.Name {
			return sprint, nil
		}
	}
	return Sprint{}, fmt.Errorf("There is no sprint '%s' found.", name)
}
