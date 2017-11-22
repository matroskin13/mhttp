package mhttp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

const Version = "0.0.1"

const configDirMacOS = "Library/Application Support/mhttp"

func getConfigFilePath() string {
	return path.Join(getConfigDir(), "mhttp_config.json")
}

func getConfigDir() string {
	currentUser, _ := user.Current()

	return path.Join(currentUser.HomeDir, configDirMacOS)
}

func GetOrCreateConfig() (*UserConfig, error) {
	config := &UserConfig{}

	configPath := getConfigFilePath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config, err = createConfig()

		if err != nil {
			return nil, err
		}
	} else {
		file, err := ioutil.ReadFile(configPath)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(file, config)

		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func createConfig() (*UserConfig, error) {
	config := &UserConfig{}

	config.Version = Version
	config.Spaces = append(config.Spaces, Space{Name: "default", Variables: make(map[string]string)})

	configBytes, err := config.toJSON()

	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(getConfigDir(), os.ModePerm)

	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(getConfigFilePath(), configBytes, 0644)

	if err != nil {
		return nil, err
	}

	return config, nil
}

type Space struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

type UserConfig struct {
	Version string  `json:"version"`
	Spaces  []Space `json:"spaces"`
}

func (c *UserConfig) GetSpace(spaceName string) *Space {
	if spaceName == "" {
		spaceName = "default"
	}

	for _, space := range c.Spaces {
		if space.Name == spaceName {
			return &space
		}
	}

	space := Space{Name: spaceName, Variables: make(map[string]string)}

	c.Spaces = append(c.Spaces, space)

	return &space
}

func (c *UserConfig) AddVar(spaceName string, name string, value string) {
	space := c.GetSpace(spaceName)

	space.Variables[name] = value
}

func (c *UserConfig) AddJSONVar(spaceName string, name string, value interface{}) error {
	jsonValue, err := json.Marshal(value)

	if err != nil {
		return err
	}

	c.AddVar(spaceName, name, string(jsonValue))

	return nil
}

func (c *UserConfig) GetVar(spaceName string, name string) (string, error) {
	space := c.GetSpace(spaceName)

	if space == nil {
		return "", errors.New("space is not defined")
	}

	return space.Variables[name], nil
}

func (c *UserConfig) Save() error {
	jsonConfig, err := c.toJSON()

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getConfigFilePath(), jsonConfig, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (c *UserConfig) toJSON() ([]byte, error) {
	return json.Marshal(c)
}
