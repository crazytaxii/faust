package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type QiniuConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	BaseUrl   string `yaml:"base_url"`
}

// Load config from yaml file.
func (fc *FaustClient) LoadConfig(file string) error {
	fc.ConfigFilePath = file
	fc.Config = &QiniuConfig{}
	// check ~/.faust/config.yaml is existed
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// create ~/.faust/config.yaml
		if err := fc.SaveConfig(); err != nil {
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
