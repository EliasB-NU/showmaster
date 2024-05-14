package config

import (
	"encoding/json"
	"log"
	"os"
)

type CFG struct {
	DB struct {
		Host      string `json:"Host"`
		Port      int    `json:"Port"`
		Username  string `json:"Username"`
		Password  string `json:"Password"`
		Daatabase string `json:"Database"`
	} `json:"DB"`

	Website struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	}

	ProjectName string `json:"Project"`
}

func GetConfig() *CFG {
	const file = "config/config.json"
	var config CFG

	cfgfile, err := os.Open(file)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error readeing config file: %d\n", err)
	}

	jsonParser := json.NewDecoder(cfgfile)
	jsonParser.Decode(&config)

	return &config
}
