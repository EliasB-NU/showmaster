package database

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	ID uint `gorm:"primaryKey;autoIncrement"`

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

	log.Println("Database initialized successfully")
	return nil
}
