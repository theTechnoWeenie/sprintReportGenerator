package jira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByProjectName(t *testing.T) {
	projects := []RapidView{RapidView{1, "Test1", true, true, true}, RapidView{2, "Test2", true, true, true}}
	projectList := AgileProjectList{projects}
	the_project, err := projectList.GetProjectByName("Test2")
	assert.Equal(t, err, nil)
	assert.Equal(t, projects[1], the_project)
}
