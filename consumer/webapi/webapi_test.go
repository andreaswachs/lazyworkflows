package webapi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/joho/godotenv"
)

// Setup and mocking

type MockRoundTripper func(r *http.Request) *http.Response

func (f MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func SetupSuite(t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		InjectHttpClient(&http.Client{
			Transport: MockRoundTripper(func(r *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("Hello, World!")),
				}
			})})
	}
}

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

func TestDispatchCanGetAWorkflow(t *testing.T) {

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
