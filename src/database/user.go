package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"showmaster/src/util"
)

// User struct for a new user
type User struct {
	Name            string
	Email           string
	Password        string
	PermissionLevel int
}

func GetUsers(db *sql.DB) ([]User, error) {
	var (
		execSQL = fmt.Sprintf(`SELECT * FROM showmaster.users;`)

		users []User
		err   error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error querying user rows: %d\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Name, &user.Email, &user.Password, &user.PermissionLevel)
		if err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Printf("Error scanning user rows: %d\n\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error iterating over user rows: %v", err)
		return nil, err
	}

	return users, nil
}

func GetUserPermissionLevel(name string, db *sql.DB) (int, error) {
	var (
		execSQL = fmt.Sprintf(`SELECT users.permlvl FROM showmaster.users WHERE name='%s';`, name)

		err error
	)
	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error querying user rows: %d\n", err)
		return 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var permlvl int
		err = rows.Scan(&permlvl)
		if err != nil {
			log.SetFlags(log.LstdFlags & log.Lshortfile)
			log.Printf("Error scanning over user rows: %d\n", err)
		}
		return permlvl, nil
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error iterating over user rows: %v", err)
		return 0, err
	}

	return 0, nil
}

// NewUser creates a user, needs the user struct from this package
func NewUser(u User, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`INSERT INTO showmaster.users (name, email, password, permlvl) VALUES ('%s', '%s', '%s', '%d')`,
			u.Name, u.Email, u.Password, u.PermissionLevel)

		err error
	)

	for _, user := range util.USERS {
		if user == u.Name {
			err = errors.New("user already exists")
			return err
		} else {
			continue
		}
	}

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error inserting user into database: %v", err)
		return err
	} else {
		util.CacheUsers(db)
		return nil
	}
}

// UpdateUser only for updating the permission level
func UpdateUser(name string, permissionLevel int, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`UPDATE showmaster.users SET permlvl = '%d' WHERE name = '%s'`, permissionLevel, name)

		err error
	)
	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error updating user into database: %v", err)
		return err
	} else {
		return nil
	}
}

func DeleteUser(name string, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`DELETE FROM showmaster.users WHERE name = '%s'`, name)

		err error
	)
	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags & log.Lshortfile)
		log.Printf("Error deleting user: %d\n", err)
		return err
	} else {
		return nil
	}
}
