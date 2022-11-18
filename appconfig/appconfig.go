package appconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/andreaswachs/lazyworkflows/meta"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"gopkg.in/yaml.v3"
)

type Repo struct {
	Token string
	Repo  string
	Owner string
}

type AppConfig struct {
	Repos []Repo
}

// Implementation of repoConfig interface on struct repo
func (r *Repo) GetToken() string {
	return r.Token
}

func (r *Repo) GetOwner() string {
	return r.Owner
}

func (r *Repo) GetRepo() string {
	return r.Repo
}

// Implementation of appConfig on struct appConfig
func (c *AppConfig) GetRepos() []Repo {
	return c.Repos
}

func (c *AppConfig) GetRepoWithRepo(name string) (Repo, error) {
	return Repo{}, errors.New("not implemented yet")
}

func (c *AppConfig) GetRepoWithOwner(owner string) (Repo, error) {
	return Repo{}, errors.New("not implemented yet")
}

func (c *AppConfig) Load() error {
	configPath := filepath.Join(xdg.DataHome, meta.AppName)
	if err := ensureCreated(configPath); err != nil {
		return err
	}

	config.AddDriver(yamlv3.Driver)
	configFilePath := filepath.Join(configPath, meta.ConfigFileName)

	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		file, err := os.Create(configFilePath)
		defer file.Close()
		if err != nil {
			return err
		}
		return err
	}

	fmt.Printf("Contents: %v\n", string(contents))

	err = yaml.Unmarshal(contents, c)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		return err
	}

	fmt.Printf("Repos: %v\n", c.Repos)

	return nil
}

func New() *AppConfig {
	return &AppConfig{}
}

// Ensures that the full path given is created or fails
func ensureCreated(configPath string) error {
	if _, err := os.Stat(configPath); err != nil {
		if err = os.MkdirAll(configPath, os.ModePerm); err != nil {
			fmt.Sprintf("Could not create configuration dir at location: \"%s\"\n", configPath)
			return err
		}
	}

	return nil
}
