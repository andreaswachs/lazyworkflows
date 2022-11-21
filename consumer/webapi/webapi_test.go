package webapi

import (
	"fmt"
	"os"
	"testing"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/joho/godotenv"
)

func TestListCanGetListOfWorkflows(t *testing.T) {
	// Test setup
	token, err := getGithubToken()
	if err != nil {
		t.Errorf("error getting Github token: %v", err)
	}

	// Test
	repo := getTestingRepo(token)
	apiConsumer := WebApi{}

	workflows, err := apiConsumer.List(repo)
	if err != nil {
		t.Errorf("error getting list of workflows: %v", err)
	}

	if len(workflows) == 0 {
		t.Errorf("error: no workflows returned")
	}
}

func TestGetCanGetAWorkflow(t *testing.T) {
	// Test setup
	token, err := getGithubToken()
	if err != nil {
		t.Errorf("error getting Github token: %v", err)
	}

	// Test
	repo := getTestingRepo(token)
	apiConsumer := WebApi{}

	workflow, err := apiConsumer.Get(repo, "33451598") // TODO: change this ID to the testing workflow should exist
	if err != nil {
		t.Errorf("error getting workflow: %v", err)
	}

	if workflow.Name == "" {
		t.Errorf("error: no workflow returned")
	}

	if workflow.Name != "Unit Tests on Push" { // Remember to change the name check after creating the testing workflow
		t.Errorf("error: wrong workflow returned: %v", workflow.Name)
	}
}

func getTestingRepo(token string) appconfig.Repo {
	return appconfig.Repo{
		Token: token,
		Repo:  "lazyworkflows",
		Owner: "andreaswachs",
	}
}

func getGithubToken() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file")
	}

	token := os.Getenv("GITHUB_TOKEN_TESTING")
	if token == "" {
		return "", fmt.Errorf("error: GITHUB_TOKEN_TESTING is not set")
	}

	return token, nil
}
