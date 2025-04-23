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
	Perms  *Permission `gorm:"not null"`
}

type BrowserToken struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`

	DeviceId string
	Token    string `gorm:"unique"`
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

	SceneID          float64
	SceneName        string
	SceneDescription string

	ExecutesIn time.Duration // The time in which this scene is loaded, starting from the selection of the previous one

	// Audio
	Audio string
	// Audio - Clients
	// Audio - Clients - Midi
	AudioMidiEnabled bool
	AudioMidiClient  uint64
	AudioMidiChannel int
	AudioMidiNote    string
	// Audio - Clients - OSC
	AudioOSCEnabled bool
	AudioOSCClient  uint64
	AudioOSCChannel int
	AudioOSCNote    string
	// Audio - Clients - GPIO
	AudioGPIOEnabled bool
	AudioGPIOClient  uint64
	AudioGPIOChannel int
	AudioGPIONote    string
	// Light
	Light string
	// Light - Clients
	// Light - Clients - Midi
	LightMidiEnabled bool
	LightMidiClient  uint64
	LightMidiChannel int
	LightMidiNote    string
	// Light - Clients - OSC
	LightOSCEnabled bool
	LightOSCClient  uint64
	LightOSCChannel int
	LightOSCNote    string
	// Light - Clients - GPIO
	LightGPIOEnabled bool
	LightGPIOClient  uint64
	LightGPIOChannel int
	LightGPIONote    string
	// Video
	Video string
	// Video - Clients
	// Video - Clients - Midi
	VideoMidiEnabled bool
	VideoMidiClient  uint64
	VideoMidiChannel int
	VideoMidiNote    string
	// Video - Clients - OSC
	VideoOSCEnabled bool
	VideoOSCClient  uint64
	VideoOSCChannel int
	VideoOSCNote    string
	// Video - Clients - GPIO
	VideoGPIOEnabled bool
	VideoGPIOClient  uint64
	VideoGPIOChannel int
	VideoGPIONote    string
	// Video - Clients - IntraCast
	VideoIntraCastEnabled bool
	VideoIntraCastClient  uint64
	VideoIntraCastChannel int
	VideoIntraCastNote    string

	EventID uint64 `gorm:"index"`
	Event   Event
}

type Client struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`

	Name  string
	Type  string
	IP    string
	Port  int
	Token string `gorm:"unique"`
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

	err = db.AutoMigrate(&Client{})
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
