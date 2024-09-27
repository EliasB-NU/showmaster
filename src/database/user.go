package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// User struct for a new user
type User struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PermissionLevel *int
}

func GetUsers(db *sql.DB) ([]User, error) {
	var (
		execSQL = fmt.Sprintf(`SELECT name, email, permlvl FROM showmaster.users;`)

		users []User
		err   error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying user rows: %d\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Name, &user.Email, &user.PermissionLevel)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning user rows: %d\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error iterating over user rows: %d\n", err)
		return nil, err
	}

	return users, nil
}

// CheckIfRegistered checks if a user is registered and if yes, the password: Password match: 202; Password mismatch: 403; Unregistered: 404
func CheckIfRegistered(email string, password string, db *sql.DB) (int, error) {
	var (
		execSQL = fmt.Sprintf(`SELECT password FROM showmaster.users WHERE email='%s';`, email)
		err     error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying user rows: %d\n", err)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var pwd string
		err = rows.Scan(&pwd)
		if err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning over user rows: %d\n", err)
			return 0, err
		}
		if password != pwd {
			return 0, nil
		} else {
			return 1, nil
		}
	}
	return 2, nil
}

// NewUser creates a user, needs the user struct from this package
func NewUser(u User, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`INSERT INTO showmaster.users (name, email, password, permlvl) VALUES ('%s', '%s', '%s', '%d')`,
			u.Name, u.Email, u.Password, u.PermissionLevel)

		err error
	)
	cacheUsers(db)

	for _, user := range users {
		if user == u.Name {
			err = errors.New("user already exists")
			return err
		} else {
			continue
		}
	}

	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error inserting user into database: %d\n", err)
		return err
	}

	return nil
}

// UpdateUser only for updating the permission level
func UpdateUser(name string, permissionLevel int, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`UPDATE showmaster.users SET permlvl = '%d' WHERE name = '%s'`, permissionLevel, name)

		err error
	)
	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error updating user into database: %d\n", err)
		return err
	}

	return nil
}

func DeleteUser(name string, db *sql.DB) error {
	var (
		execSQL = fmt.Sprintf(`DELETE FROM showmaster.users WHERE name = '%s'`, name)

		err error
	)
	_, err = db.Exec(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error deleting user: %d\n", err)
		return err
	}

	return nil
}

// USERS list of all users currently registered
var users []string

// CacheUsers caches the users into a var stored in this file, will be called after every user action
func cacheUsers(db *sql.DB) {
	var (
		execSQL = fmt.Sprintf(`SELECT name FROM showmaster.users`)

		err error
	)

	rows, err := db.Query(execSQL)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error querying rows from showmaster.users: %d\n", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Processing the results
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Printf("Error scanning row: %d\n", err)
		}
		users = append(users, name)
	}

	if err := rows.Err(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Error with rows: %d\n", err)
	}
}
