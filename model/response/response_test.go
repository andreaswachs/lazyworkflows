package response

import (
	"testing"
)

func TestCanDeserializeGetResponse(t *testing.T) {
	responseText := getWorkflow1()

	var responseObj Get
	FromString(responseText, &responseObj)

	if responseObj.Workflow.Id != "161335" {
		t.Fatalf("Expected workflow id to be 161335, but got %v", responseObj.Workflow.Id)
	}
	if responseObj.Workflow.Name != "CI" {
		t.Fatalf("Expected workflow name to be CI, but got %v", responseObj.Workflow.Name)
	}
	if responseObj.Workflow.Path != ".github/workflows/blank.yaml" {
		t.Fatalf("Expected workflow path to be .github/workflows/blank.yaml, but got %v", responseObj.Workflow.Path)
	}
	if responseObj.Workflow.State != "active" {
		t.Fatalf("Expected workflow state to be active, but got %v", responseObj.Workflow.State)
	}
	if responseObj.Workflow.CreatedAt != "2020-01-08T23:48:37.000-08:00" {
		t.Fatalf("Expected workflow created_at to be 2020-01-08T23:48:37.000-08:00, but got %v", responseObj.Workflow.CreatedAt)
	}
	if responseObj.Workflow.UpdatedAt != "2020-01-08T23:50:21.000-08:00" {
		t.Fatalf("Expected workflow updated_at to be 2020-01-08T23:50:21.000-08:00, but got %v", responseObj.Workflow.UpdatedAt)
	}
	if responseObj.Workflow.Url != "https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335" {
		t.Fatalf("Expected workflow url to be https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335, but got %v", responseObj.Workflow.Url)
	}
	if responseObj.Workflow.HtmlUrl != "https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335" {
		t.Fatalf("Expected workflow html_url to be https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335, but got %v", responseObj.Workflow.HtmlUrl)
	}
	if responseObj.Workflow.BadgeUrl != "https://github.com/octo-org/octo-repo/workflows/CI/badge.svg" {
		t.Fatalf("Expected workflow badge_url to be https://github.com/octo-org/octo-repo/workflows/CI/badge.svg, but got %v", responseObj.Workflow.BadgeUrl)
	}
}

func TestCanDeserializeListResponse(t *testing.T) {
	responseText := listWorkflows()

	var responseObj List
	FromString(responseText, &responseObj)

	if responseObj.TotalCount != 1 {
		t.Fatalf("Expected total count to be 2, but got %v", responseObj.TotalCount)
	}
	if len(responseObj.Workflows) != 1 {
		t.Fatalf("Expected 2 workflows, but got %v", len(responseObj.Workflows))
	}
	// TODO: Can we make assume that the flowflow is deserialized correctly?
}

func TestCanDeserializeWorkflowDispatchResponse(t *testing.T) {
	responseText := getStatus200Response()

	var responseObj Dispatch

	FromString(responseText, &responseObj)

	if responseObj.Status != 200 {
		t.Fatalf("Expected status to be 200, but got %v", responseObj.Status)
	}
}

func TestCanDeserializeWorkflowEnableResponse(t *testing.T) {
	responseText := getStatus200Response()

	var responseObj Enable

	FromString(responseText, &responseObj)

	if responseObj.Status != 200 {
		t.Fatalf("Expected status to be 200, but got %v", responseObj.Status)
	}
}

func TestCanDeserializeWorkflowDisableResponse(t *testing.T) {
	responseText := getStatus200Response()

	var responseObj Disable

	FromString(responseText, &responseObj)

	if responseObj.Status != 200 {
		t.Fatalf("Expected status to be 200, but got %v", responseObj.Status)
	}
}

// Actual response text from GitHub API
//
func getWorkflow1() string {
	return `{"id":161335,"node_id":"MDg6V29ya2Zsb3cxNjEzMzU=","name":"CI","path":".github/workflows/blank.yaml","state":"active","created_at":"2020-01-08T23:48:37.000-08:00","updated_at":"2020-01-08T23:50:21.000-08:00","url":"https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335","html_url":"https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335","badge_url":"https://github.com/octo-org/octo-repo/workflows/CI/badge.svg"}`
}

func getWorkflow2() string {
	return `{"id":20,"node_id":"MDg6V29ya2Zsb3cxNjEzMzU=","name":"CD","path":".github/workflows/other.yaml","state":"disabled","created_at":"2020-01-08T23:48:37.000-08:00","updated_at":"2020-01-08T23:50:21.000-08:00","url":"https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335","html_url":"https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335","badge_url":"https://github.com/octo-org/octo-repo/workflows/CI/badge.svg"}`
}

func getStatus200Response() string {
	return `{"status": 200}`
}

func getRealListWorkflows() string {
	return "{\"total_count\":1,\"workflows\":[{\"id\":33451598,\"node_id\":\"W_kwDOH0TxRs4B_m5O\",\"name\":\"Unit Tests on Push\",\"path\":\".github/workflows/unit-tests-on-push.yml\",\"state\":\"active\",\"created_at\":\"2022-08-28T08:55:14.000+02:00\",\"updated_at\":\"2022-08-28T10:26:51.000+02:00\",\"url\":\"https://api.github.com/repos/andreaswachs/lazyworkflows/actions/workflows/33451598\",\"html_url\":\"https://github.com/andreaswachs/lazyworkflows/blob/main/.github/workflows/unit-tests-on-push.yml\",\"badge_url\":\"https://github.com/andreaswachs/lazyworkflows/workflows/Unit%20Tests%20on%20Push/badge.svg\"}]}"
}

func listWorkflows() string {
	return getRealListWorkflows()
}
