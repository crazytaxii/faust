package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

func loadConfig(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		err = createConfig(fileName)
		if err != nil {
			return err
		}
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, conf)
}

func createConfig(fileName string) error {
	conf := &Config{}
	j, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, j, 0644)
}

func saveConfig(c *cli.Context) error {
	j, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, j, 0644)
}
