package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/howeyc/gopass"
	"github.com/theTechnoWeenie/sprintReportGenerator/jira"
)

func main() {
	//Command line options
	projectNamePtr := flag.String("project", "", "The name of the agile project (not the TLA)")
	sprintNamePtr := flag.String("sprint", "", "The name of the sprint to retreive issues for")
	jiraUrlPtr := flag.String("jiraUrl", "", "The url of the jira instance to hit.")
	version := flag.Bool("version", false, "Displays the version")
	help := flag.Bool("help", false, "Displays help.")

	flag.Parse()
	//End command line options
	if *version {
		fmt.Println("Version: ", VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "Example: %s -jiraUrl 'https://my.jira.com/' -project 'My Team' -sprint 'Our Last Sprint'\n", os.Args[0])
		os.Exit(0)
	}

	variables := map[string]*string{
		"No project provided":     projectNamePtr,
		"No sprint name provided": sprintNamePtr,
		"No jira url provided":    jiraUrlPtr,
	}

	for key, val := range variables {
		if *val == "" {
			fmt.Println(key)
			fmt.Println()
			os.Exit(1)
		}
	}

	user, pass := getUserAndPassword()
	credentials := jira.Credentials{}
	credentials.Init(user, pass)
	url := *jiraUrlPtr
	if url[len(url)-1:] != "/" {
		url += "/"
	}
	client := jira.NewClient(url, credentials)

	fmt.Println("Populating Jira projects...")
	projects, err := getJiraProjects(client)
	exitIfError(err, "Could not get project!")
	fmt.Printf("Found %d projects\n", len(projects.Projects))
	desired_project, err := projects.GetProjectByName(*projectNamePtr)
	exitIfError(err, fmt.Sprintf("Project '%s' is not found", *projectNamePtr))

	fmt.Printf("retrieved %s\n", desired_project.Name)
	fmt.Println("Finding sprint ", *sprintNamePtr)
	sprints := getSprintsForRapidView(desired_project, client)
	desired_sprint, err := sprints.FindSprintWithName(*sprintNamePtr)
	exitIfError(err, fmt.Sprintf("Sprint '%s' is not found", *sprintNamePtr))

	fmt.Println("found.\nGetting report details...")
	sprintReportQuery := map[string]string{"rapidViewId": fmt.Sprintf("%d", desired_project.Id), "sprintId": fmt.Sprintf("%d", desired_sprint.Id)}
	reportJson := client.Get("rest/greenhopper/1.0/rapid/charts/sprintreport", sprintReportQuery)
	fmt.Println("done.\nGetting velocity data...")
	velocity := jira.ParseVelocity(client.Get("rest/greenhopper/1.0/rapid/charts/velocity", sprintReportQuery))
	fmt.Println("done.")
	var report jira.Report
	err = json.Unmarshal(reportJson, &report)
	exitIfError(err, "Could not get sprint report.")
	fmt.Print("\n\nWhile editing a confluence page, press <ctrl>+<shift>+d and paste the following into the box:\n\n")
	fmt.Println(report.GenerateReport(client, velocity))
	if reportJson == nil {
		log.Fatal("No report returned")
	}
}

func getJiraProjects(client jira.Client) (jira.AgileProjectList, error) {
	raw_json := client.Get("rest/greenhopper/1.0/rapidview", nil)
	var agileViews jira.AgileProjectList
	err := json.Unmarshal(raw_json, &agileViews)
	return agileViews, err
}

func getSprintsForRapidView(desired_project jira.RapidView, client jira.Client) jira.SprintList {
	sprintsJson := client.Get(fmt.Sprintf("rest/greenhopper/1.0/sprintquery/%d", desired_project.Id), nil)
	var sprints jira.SprintList
	if err := json.Unmarshal(sprintsJson, &sprints); err != nil {
		log.Fatal(err)
	}
	return sprints
}

func getUserAndPassword() (username string, password string) {
	fmt.Printf("User: ")
	user, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Printf("Password: ")
	pass := gopass.GetPasswdMasked()
	return user, string(pass)
}

func exitIfError(err error, message string) {
	if err != nil {
		if message != "" {
			log.Fatal(fmt.Sprintf("%s\nCaused by: %v", message, err))
		}
		log.Fatal(err)
	}
}
