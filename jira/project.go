package jira

import "fmt"

type RapidView struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Editable      bool   `json:"canEdit"`
	SprintSupport bool   `json:"sprintSupportEnabled"`
	ShowDays      bool   `json:"showDaysInColumn"`
}

type AgileProjectList struct {
	Projects []RapidView `json:"views"`
}

func (projectList *AgileProjectList) GetProjectByName(name string) (project RapidView, err error) {
	for _, project := range projectList.Projects {
		if project.Name == name {
			return project, nil
		}
	}

	return RapidView{}, fmt.Errorf("Project '%s' was not found.", name)
}
