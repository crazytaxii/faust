package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Load config from yaml file.
func (fc *FaustClient) LoadConfig(file string) error {
	fc.ConfigFilePath = file
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		err = fc.SaveConfig()
		if err != nil {
			return err
		}
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, &fc.Config)
}

// Create a new config file.
func (fc *FaustClient) SaveConfig() error {
	y, err := yaml.Marshal(fc.Config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fc.ConfigFilePath, y, 0644)
}
