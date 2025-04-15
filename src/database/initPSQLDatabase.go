package database

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	Username string
	Email    string `gorm:"unique"`
	Password string

	Tokens *[]BrowserToken
	Perms  *Permission
}

type BrowserToken struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`

	DeviceId string
	Key      string
	UserID   uint64 `gorm:"index"`
	User     User
}

type Permission struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`

	Login  bool
	Admin  bool
	Events int

	UserID uint64 `gorm:"index"`
	User   User
}

type Event struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`

	Name        string
	Description string
	Location    string
	StartDate   time.Time
	EndDate     time.Time

	CurrentScene int
	TimeElapsed  time.Duration

	Scenes []SceneEntries
}

type SceneEntries struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`

	SceneName        string
	SceneDescription string

	EventID uint64 `gorm:"index"`
	Event   Event
}

func InitPSQLDatabase(db *gorm.DB) error {
	var err error

	err = db.AutoMigrate(&User{})
	if err != nil {
		return errors.New("failed to auto migrate users table: " + err.Error())
	}

	err = db.AutoMigrate(&BrowserToken{})
	if err != nil {
		return errors.New("failed to auto migrate browser token table: " + err.Error())
	}

	err = db.AutoMigrate(&Permission{})
	if err != nil {
		return errors.New("failed to auto migrate permissions table: " + err.Error())
	}

	err = db.AutoMigrate(&Event{})
	if err != nil {
		return errors.New("failed to auto migrate events table: " + err.Error())
	}

	err = db.AutoMigrate(&SceneEntries{})
	if err != nil {
		return errors.New("failed to auto migrate scenes table: " + err.Error())
	}

	// Create initial admin user, if not exists (email: admin@example.com, username: admin, password: admin)
	result := db.Where("username = ?", "admin").First(&User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var user User
		user.Username = "admin"
		user.Password, _ = hashPassword("admin")
		user.Email = "admin@example.com"
		var perms Permission
		perms.Login = true
		perms.Admin = true
		perms.Events = 3
		user.Perms = &perms
		err = db.Create(&user).Error
		if err != nil {
			return errors.New("failed to create initial admin user: " + err.Error())
		}
	} else if result.Error != nil {
		return errors.New("failed to check for initial admin user: " + result.Error.Error())
	} else {
		log.Println("Initial admin user already exists")
	}

	log.Println("Database initialized successfully")
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
