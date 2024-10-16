package util

import (
	"database/sql"
	"fmt"
	"log"
)

func GetAdminUsers(db *sql.DB) map[string]string {
	var (
		adminUsers = make(map[string]string)

		execSQL = fmt.Sprintf(`SELECT email, password FROM showmaster.users WHERE permlvl = 3`)

		err error
	)
	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows: %d\n", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error closing rows: %d\n", err)
		}
	}(rows)

	for rows.Next() {
		var username string
		var password string
		err = rows.Scan(&username, &password)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning rows: %d\n", err)
		}
		adminUsers[username] = password
	}
	return adminUsers
}
