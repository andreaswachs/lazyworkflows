package response

import (
	"testing"

	"github.com/andreaswachs/lazyworkflows/test_resources"
)

func TestCanDeserializeGetResponse(t *testing.T) {
	responseText := test_resources.Workflow1

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
	responseText := test_resources.ListResponse

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
	responseText := test_resources.Status200Response

	var responseObj Dispatch

	FromString(responseText, &responseObj)

	if responseObj.Status != 200 {
		t.Fatalf("Expected status to be 200, but got %v", responseObj.Status)
	}
}

func TestCanDeserializeWorkflowEnableResponse(t *testing.T) {
	responseText := test_resources.Status200Response

	var responseObj Enable

	FromString(responseText, &responseObj)

	if responseObj.Status != 200 {
		t.Fatalf("Expected status to be 200, but got %v", responseObj.Status)
	}
}

func TestCanDeserializeWorkflowDisableResponse(t *testing.T) {
	responseText := test_resources.Status200Response

	var responseObj Disable

	FromString(responseText, &responseObj)

	if responseObj.Status != 200 {
		t.Fatalf("Expected status to be 200, but got %v", responseObj.Status)
	}
}
