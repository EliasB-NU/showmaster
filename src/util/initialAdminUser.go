package util

import (
	"database/sql"
	"log"
	"showmaster/src/config"
	"showmaster/src/database"
)

func CreateInitialAdminUser(db *sql.DB, cfg *config.CFG) {
	var (
		permlvl = 3
		u       = database.User{
			Name:            cfg.User.AdminUserName,
			Email:           "admin@example.com",
			Password:        cfg.User.AdminPassword,
			PermissionLevel: &permlvl,
		}

		err error
	)
	err = database.NewUser(u, db)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error creating admin user: %d\n", err)
	}
}
