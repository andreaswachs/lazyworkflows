package webapi

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/model/response"
	"github.com/andreaswachs/lazyworkflows/test_resources"
)

// Setup and mocking

type MockRoundTripper func(r *http.Request) *http.Response

func (f MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

// Mocks the sharedHttpClient such that we can control the response
func SetupSuite(t *testing.T, response string) func(t *testing.T) {
	InjectHttpClient(&http.Client{
		Transport: MockRoundTripper(func(r *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(response)),
			}
		})})

	return func(t *testing.T) {}
}

func TestListCanGetListOfWorkflows(t *testing.T) {

	responseInterface, err := executeWithSetup(t, test_resources.ListResponse, func(apiConsumer WebApi, repo appconfig.Repo) (interface{}, error) {
		return apiConsumer.List(repo)
	})
	if err != nil {
		t.Errorf("error getting list of workflows: %v", err)
	}

	workflows := responseInterface.([]response.Workflow)
	if len(workflows) == 0 {
		t.Errorf("error: no workflows returned")
	}
}

func TestGetCanGetAWorkflow(t *testing.T) {
	responseInterface, err := executeWithSetup(t, test_resources.Workflow1, func(apiConsumer WebApi, repo appconfig.Repo) (interface{}, error) {
		return apiConsumer.Get(repo, "filler")
	})
	if err != nil {
		t.Errorf("error getting workflow: %v", err)
	}

	if err != nil {
		t.Errorf("error getting workflow: %v", err)
	}

	workflow := responseInterface.(response.Workflow)
	if workflow.Name == "" {
		t.Errorf("error: no workflow returned")
	}
	if workflow.Name != "CI" { // Remember to change the name check after creating the testing workflow
		t.Errorf("error: expected name \"CI\", got: %v", workflow.Name)
	}
}

func TestDispatchCanGetAWorkflow(t *testing.T) {
	responseInterface, err := executeWithSetup(t, test_resources.Status200Response, func(apiConsumer WebApi, repo appconfig.Repo) (interface{}, error) {
		return apiConsumer.Dispatch(repo, "filler")
	})
	if err != nil {
		t.Errorf("error dispatching workflow: %v", err)
	}

	dispatchResponse := responseInterface.(response.Dispatch)
	if dispatchResponse.Status != 200 {
		t.Errorf("error: expected status 200, got: %v", dispatchResponse.Status)
	}
}

func TestEnableCanEnableAWorkflow(t *testing.T) {
	responseInterface, err := executeWithSetup(t, test_resources.Status200Response, func(apiConsumer WebApi, repo appconfig.Repo) (interface{}, error) {
		return apiConsumer.Enable(repo, "filler")
	})
	if err != nil {
		t.Errorf("error enabling workflow: %v", err)
	}

	enableResponse := responseInterface.(response.Enable)
	if enableResponse.Status != 200 {
		t.Errorf("error: expected status 200, got: %v", enableResponse.Status)
	}
}

func TestDisableCanDisableAWorkflow(t *testing.T) {
	responseInterface, err := executeWithSetup(t, test_resources.Status200Response, func(apiConsumer WebApi, repo appconfig.Repo) (interface{}, error) {
		return apiConsumer.Disable(repo, "filler")
	})
	if err != nil {
		t.Errorf("error disabling workflow: %v", err)
	}

	disableResponse := responseInterface.(response.Disable)
	if disableResponse.Status != 200 {
		t.Errorf("error: expected status 200, got: %v", disableResponse.Status)
	}
}

func executeWithSetup(t *testing.T, requestResponse string, f func(apiConsumer WebApi, repo appconfig.Repo) (interface{}, error)) (interface{}, error) {
	teardown := SetupSuite(t, requestResponse)
	defer teardown(t)

	return f(WebApi{}, getTestingRepo())
}

func getTestingRepo() appconfig.Repo {
	return appconfig.Repo{
		Token: "filler",
		Repo:  "filler",
		Owner: "filler",
	}
}
