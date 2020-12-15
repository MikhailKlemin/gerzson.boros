package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//GeneralConfig is configuration
type GeneralConfig struct {
	OutDir             string `json:"output_directory"`
	DomainPath         string `json:"path_to_domain_file"`
	DatabaseName       string `json:"database_name"`
	DatabaseHost       string `json:"database_host"`
	DatabaseCollection string `json:"database_collection"`
	LevelDBPath        string `json:"level_db_path"`
	Concurrency        int    `json:"concurrency"`
	ChromeTimeout      int    `json:"chrome_timeout"`
}

//LoadGeneralConfig loads gen config
func LoadGeneralConfig() GeneralConfig {
	var c GeneralConfig
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Cannot read config file " + err.Error())
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal("Cannot read json file " + err.Error())

	}

	return c
}
