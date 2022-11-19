package appconfig

import (
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

func (c *AppConfig) Load() error {
	// Do some initial configuration to enable reading the config file
	config.AddDriver(yamlv3.Driver)

	// Ensures that the directory path before the config file exists
	configPath := filepath.Join(xdg.DataHome, meta.AppName)
	if err := ensureCreated(configPath); err != nil {
		return err
	}

	// Cread the config file or create it. If it didn't exist, it will be created
	// and error out to encourage the user to fill it out
	configFilePath := filepath.Join(configPath, meta.ConfigFileName)
	contents, err := readConfig(configFilePath)
	if err != nil {
		return err
	}

	// Unmarshal/deserialize the config file, turning it into our app config struct
	if err = yaml.Unmarshal(contents, c); err != nil {
		return fmt.Errorf("while reading the config.yml file, an error occured: %w", err)
	}

	return nil
}

func New() *AppConfig {
	return &AppConfig{}
}

// Reads the config file or creates it if it doesn't exist,
// given the path to the config file
func readConfig(configFilePath string) ([]byte, error) {
	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("could not read config file: %w", err)
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("could not create config file: %w", err)
	}

	file.Close()
	return contents, nil
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
