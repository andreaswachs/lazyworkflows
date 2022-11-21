package consumer

import (
	"github.com/andreaswachs/lazyworkflows/appconfig"
	"github.com/andreaswachs/lazyworkflows/consumer/webapi"
	"github.com/andreaswachs/lazyworkflows/model/response"
)

type Consumer interface {
	List(appconfig.Repo) ([]response.Workflow, error)
	Get(appconfig.Repo, string) (response.Workflow, error)
	Dispatch(appconfig.Repo, string) (response.Dispatch, error)
	Enable(appconfig.Repo, string) (response.Enable, error)
	Disable(appconfig.Repo, string) (response.Disable, error)
}

// Returns a new API consumer
// The concrete consumer can be configured in this functiok
func New() Consumer {
	return &webapi.WebApi{}
}
