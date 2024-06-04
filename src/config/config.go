package config

import (
	"encoding/json"
	"log"
	"os"
)

type CFG struct {
	DB struct {
		Host     string `json:"Host"`
		Port     int    `json:"Port"`
		Username string `json:"Username"`
		Password string `json:"Password"`
		Database string `json:"Database"`
	} `json:"DB"`

	Website struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	}

	ProjectName string `json:"Project"`
}

// func GetConfig() *CFG {
// 	var c CFG
// 	c.DB.Host = "db"
// 	c.DB.Port = 5432
// 	c.DB.Username = os.Getenv("DBUser")
// 	c.DB.Password = os.Getenv("DBPassword")
// 	c.DB.Database = os.Getenv("Database")
// 	c.Website.Host = "0.0.0.0"
// 	c.Website.Port = 80
// 	c.ProjectName = os.Getenv("ProjectName")

// 	return &c
// }

// Enable for local testing without docker container
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
