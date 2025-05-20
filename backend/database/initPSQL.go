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
	Email    string
	Password string

	Tokens *[]BrowserToken

	Permission *Permission
}

type BrowserToken struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	DeviceId string
	Token    string `gorm:"unique"`
	UserID   uint64 `gorm:"index"`
	User     User
}

// Permission
// There are 4 different roles:
// Admin: Can do everything and edit/delete/create events, users & clients
// Director: Can work on events and select new scenes
// FollowOperator: Can see the follow operator plans
// Viewer: Can just see all the scenes with their basic info
type Permission struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	Role  string
	Level int

	UserID uint64 `gorm:"index"`
	User   User
}

type Event struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	Name        string `json:"name"`
	Description string `json:"description"`

	StartDate string        `json:"start_date"`
	EndDate   string        `json:"end_date"`
	Timer     time.Duration `json:"timer"`

	Scenes []Scene `json:"scenes"`
}

type Scene struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Audio string `json:"audio"`
	Light string `json:"light"`
	Video string `json:"video"`

	MIDI        []MIDI       `json:"midi"`
	OSC         []OSC        `json:"osc"`
	GPIO        []GPIO       `json:"gpio"`
	FollowSpots []FollowSpot `json:"follow_spots"`

	EventID uint64 `gorm:"index"`
	Event   Event
}

type MIDI struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SceneID uint64 `gorm:"index"`
	Scene   Scene  `json:"scene"`
}

type OSC struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SceneID uint64 `gorm:"index"`
	Scene   Scene  `json:"scene"`
}

type GPIO struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SceneID uint64 `gorm:"index"`
	Scene   Scene  `json:"scene"`
}

type FollowSpot struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	SceneID uint64 `gorm:"index"`
	Scene   Scene  `json:"scene"`
}

type Clients struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement"`
}

// InitPSQLDatabase Creates all the necessary tables for the app to work
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
		return errors.New("failed to auto migrate permission table: " + err.Error())
	}

	err = db.AutoMigrate(&Event{})
	if err != nil {
		return errors.New("failed to auto migrate event table: " + err.Error())
	}

	err = db.AutoMigrate(&Scene{})
	if err != nil {
		return errors.New("failed to auto migrate scene table: " + err.Error())
	}

	err = db.AutoMigrate(&MIDI{})
	if err != nil {
		return errors.New("failed to auto migrate midi table: " + err.Error())
	}

	err = db.AutoMigrate(&OSC{})
	if err != nil {
		return errors.New("failed to auto migrate osc table: " + err.Error())
	}

	err = db.AutoMigrate(&GPIO{})
	if err != nil {
		return errors.New("failed to auto migrate gpio table: " + err.Error())
	}

	err = db.AutoMigrate(&Clients{})
	if err != nil {
		return errors.New("failed to auto migrate clients table: " + err.Error())
	}

	// Create initial admin user, if not exists (email: admin@example.com, username: admin, password: admin)
	result := db.Where("username = ?", "admin").First(&User{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		var user User
		user.Username = "admin"
		user.Password, _ = hashPassword("admin")
		user.Email = "admin@example.com"
		var perms = Permission{
			Role:  "Admin",
			Level: 0,
		}
		user.Permission = &perms
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
