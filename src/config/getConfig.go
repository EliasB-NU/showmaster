package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		PSQL struct {
			Host     string `yaml:"Host"`
			Port     int    `yaml:"Port"`
			User     string `yaml:"User"`
			Password string `yaml:"Password"`
			DBName   string `yaml:"DBName"`
			TimeZone string `yaml:"Timezone"`
		} `yaml:"PSQL"`
		Redis struct {
			Host     string `yaml:"Host"`
			Port     string `yaml:"Port"`
			DB       int    `yaml:"DB"`
			User     string `yaml:"User"`
			Password string `yaml:"Password"`
		} `yaml:"Redis"`
	} `yaml:"Database"`
}

func GetConfig() *Config {
	var config *Config

	if os.Args[1] == "dev" {
		file, err := os.Open("./config/config.yaml")
		if err != nil {
			log.Fatalf("Error opening config file: %v", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Error closing config file: %v", err)
			}
		}(file)

		// Read the file and parse the YAML content
		yamlParser := yaml.NewDecoder(file)
		err = yamlParser.Decode(&config)
		if err != nil {
			log.Fatalf("Error parsing YAML file: %v", err)
		}
	} else if os.Args[1] == "prod" {
		config.Database.PSQL.Host = os.Getenv("PSQL_HOST")
		config.Database.PSQL.Port, _ = strconv.Atoi(os.Getenv("PSQL_PORT"))
		config.Database.PSQL.User = os.Getenv("PSQL_USER")
		config.Database.PSQL.Password = os.Getenv("PSQL_PASSWORD")
		config.Database.PSQL.DBName = os.Getenv("PSQL_DBNAME")
		config.Database.PSQL.TimeZone = os.Getenv("PSQL_TIMEZONE")
		config.Database.Redis.Host = os.Getenv("REDIS_HOST")
		config.Database.Redis.Port = os.Getenv("REDIS_PORT")
		config.Database.Redis.DB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
		config.Database.Redis.User = os.Getenv("REDIS_USER")
		config.Database.Redis.Password = os.Getenv("REDIS_PASSWORD")
	} else {
		panic("Invalid environment")
		return nil
	}

	return config
}
