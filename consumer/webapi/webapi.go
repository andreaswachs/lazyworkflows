package webapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/model/response"
)

type action uint8

// The http client is a global variable and thus able to get mocked by tests
var sharedHttpClient *http.Client

const (
	enable action = iota
	disable
	dispatch
	get
	list
)

// The data structure for the WebApi consumer.
// This implements the Consumer interface
type WebApi struct {
}

type webApiRequest struct {
	Repo appconfig.Repo
	Id   string
}

// List returns a list of workflows for a given repo
func (w *WebApi) List(repo appconfig.Repo) ([]response.Workflow, error) {
	apiResponse, err := doRequest(list, repo, "")
	if err != nil {
		return nil, err
	}

	lstResponse := response.List{}
	err = response.FromString(apiResponse, &lstResponse)
	if err != nil {
		return nil, err
	}

	return lstResponse.Workflows, nil
}

// Get returns a single workflow for a given repo
func (w *WebApi) Get(repo appconfig.Repo, id string) (response.Workflow, error) {
	apiResponse, err := doRequest(get, repo, id)
	if err != nil {
		return response.Workflow{}, err
	}

	getResponse := response.Get{}
	err = response.FromString(apiResponse, &getResponse)
	if err != nil {
		return response.Workflow{}, err
	}

	return getResponse.Workflow, nil
}

// Dispatch triggers a workflow for a given repo
func (w *WebApi) Dispatch(repo appconfig.Repo, id string) (response.Dispatch, error) {
	return response.Dispatch{}, nil
}

// Enable enables a workflow for a given repo
func (w *WebApi) Enable(repo appconfig.Repo, id string) (response.Enable, error) {
	return response.Enable{}, nil
}

// Disable disables a workflow for a given repo
func (w *WebApi) Disable(repo appconfig.Repo, id string) (response.Disable, error) {
	return response.Disable{}, nil
}

func GetHttpClient() *http.Client {
	if sharedHttpClient == nil {
		sharedHttpClient = http.DefaultClient
	}
	return sharedHttpClient
}

func InjectHttpClient(injectedClient *http.Client) {
	sharedHttpClient = injectedClient
}

func doRequest(target action, repo appconfig.Repo, id string) (string, error) {
	url, err := newWebApiRequest().
		withRepo(repo).
		withId(id).
		build(target)

	if err != nil {
		return "", err
	}

	var method string
	var bodyRaw []byte

	switch target {
	case disable:
	case enable:
		method = "PUT"
	case dispatch:
		method = "POST"
	case get:
	case list:
		method = "GET"
	default:
		return "", fmt.Errorf("invalid target")
	}

	body := bytes.NewReader(bodyRaw)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repo.Token))

	resp, err := sharedHttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseText), nil
}

// Use the builder pattern to create a new webApiRequest
func newWebApiRequest() *webApiRequest {
	return &webApiRequest{}
}

// Set the repo for the webApiRequest
func (w *webApiRequest) withRepo(repo appconfig.Repo) *webApiRequest {
	w.Repo = repo
	return w
}

// Set the id for the webApiRequest
func (w *webApiRequest) withId(id string) *webApiRequest {
	w.Id = id
	return w
}

// Build the webApiRequest
func (w *webApiRequest) build(target action) (string, error) {
	// Check to see if the repo is set and valid (not empty)
	err := checkValidRepo(w.Repo)
	if err != nil {
		return "", err
	}

	switch target {
	case enable:
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/enable", w.Repo.Owner, w.Repo.Repo, w.Id), nil
	case disable:
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/disable", w.Repo.Owner, w.Repo.Repo, w.Id), nil
	case dispatch:
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/dispatches", w.Repo.Owner, w.Repo.Repo, w.Id), nil
	case get:
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s", w.Repo.Owner, w.Repo.Repo, w.Id), nil
	case list:
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows", w.Repo.Owner, w.Repo.Repo), nil
	default:
		return "", fmt.Errorf("invalid target")
	}
}

func checkValidRepo(repo appconfig.Repo) error {
	if repo.Owner == "" {
		return fmt.Errorf("owner is not set for repository settings: %v", repo)
	}
	if repo.Repo == "" {
		return fmt.Errorf("repo is not set for repository settings: %v", repo)
	}
	if repo.Token == "" {
		return fmt.Errorf("token is not set for repository settings: %v", repo)
	}
	return nil
}
