package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
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

var config CFG = *GetConfig()

func InitDB() *sql.DB {
	// Connect to PostgreSQL database
	psql := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatalf("Error opening Database: %v\n", err)
	}

	return db
}

func main() {
	db := InitDB()
	for i := 1; i < 500; i++ {
		sql := fmt.Sprintf("INSERT INTO %s (id, name, audio, licht, pptx, notes) VALUES (%d, 'the name', 'your audio changes', 'maybe some lightning', 'do you have video stuff?', 'maybe someone has to move some probs');", "secondaryproject", i)
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}
	defer db.Close()
}
