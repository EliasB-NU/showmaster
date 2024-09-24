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
		TimeZone string `json:"TimeZone"`
	} `json:"DB"`

	User struct {
		AdminUserName string `json:"AdminUserName"`
		AdminPassword string `json:"AdminPassword"`
	} `json:"User"`

	Website struct {
		Host string `json:"Host"`
		Port int    `json:"Port"`
	}
}

func GetConfig() *CFG {
	if os.Args[1] == "prod" {
		var c CFG
		c.DB.Host = "db"
		c.DB.Port = 5432
		c.DB.Username = os.Getenv("DBUser")
		c.DB.Password = os.Getenv("DBPassword")
		c.DB.Database = os.Getenv("Database")
		c.DB.TimeZone = os.Getenv("TimeZone")
		c.User.AdminUserName = os.Getenv("AdminUserName")
		c.User.AdminPassword = os.Getenv("AdminPassword")
		c.Website.Host = "0.0.0.0"
		c.Website.Port = 80

		return &c
	} else if os.Args[1] == "dev" {
		const file = "config/config.json"
		var config CFG

		cfgfile, err := os.Open(file)
		if err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Fatalf("Error readeing config file: %d\n", err)
		}

		jsonParser := json.NewDecoder(cfgfile)
		err = jsonParser.Decode(&config)
		if err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Fatalf("Error readeing config file: %d\n", err)
		}

		return &config
	} else {
		panic("Error: Wrong command line argument")
		return nil
	}
}
