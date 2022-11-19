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
		fmt.Fprintf(os.Stderr, "Could not create config dir: %v", err)
		return err
	}

	// Cread the config file or create it. If it didn't exist, it will be created
	// and error out to encourage the user to fill it out
	configFilePath := filepath.Join(configPath, meta.ConfigFileName)
	contents, err := readConfig(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read config file: %v", err)
		return err
	}

	// Unmarshal/deserialize the config file, turning it into our app config struct
	if err = yaml.Unmarshal(contents, c); err != nil {
		fmt.Fprintf(os.Stderr, "while reading the %s file, an error occurred: %v", configFilePath, err)
		return err
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
		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create config file: %v", err)
			return make([]byte, 0), err
		}

		file.Close()
		fmt.Fprintf(os.Stderr, "Config file created at %s. Please fill it out and try again. System error: %v", configFilePath, err)
		return make([]byte, 0), err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create config file: %v", err)
		return make([]byte, 0), err
	}

	file.Close()
	return contents, nil
}

// Ensures that the full path given is created or fails
func ensureCreated(configPath string) error {
	if _, err := os.Stat(configPath); err != nil {
		if err = os.MkdirAll(configPath, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Could not create config dir at location: %s. Full error: %v", configPath, err)
			return err
		}
		fmt.Fprintf(os.Stderr, "could not query file system for existence of config dir: %s. Full error message: %v", configPath, err)
		return err
	}

	return nil
}
