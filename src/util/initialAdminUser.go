package util

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"showmaster/src/config"
)

func CreateInitialAdminUser(db *sql.DB, cfg *config.CFG) error {
	var (
		execSQL = fmt.Sprintf(`INSERT INTO showmaster.users (name, email, password, permlvl) VALUES ('%s', 'admin@example.com', '%s', '3');`,
			cfg.User.AdminUserName, cfg.User.AdminPassword)

		err error
	)

	for _, user := range USERS {
		if user == cfg.User.AdminUserName {
			err = errors.New("user already exists")
			return err
		} else {
			continue
		}
	}

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating initial admin user: %d\n", err)
		return err
	}

	CacheUsers(db)
	return nil
}
