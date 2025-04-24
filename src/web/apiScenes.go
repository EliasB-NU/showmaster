package web

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"showmaster/src/database"
	"showmaster/src/util"
	"strconv"
	"time"
)

func (a *API) getScenes(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 1, util.Events, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	// Get all scenes based on the id parameter
	var scenes []database.SceneEntrie
	err := a.DB.Where("event_id = ?", c.Params("id")).Find(&scenes).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Error getting scenes: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	// Get event for active scene
	var event database.Event
	err = a.DB.Find(&event, c.Params("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Error getting event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"scenes":      scenes,
		"activeScene": event.CurrentScene,
	})
}

func (a *API) createScene(c *fiber.Ctx) error {
	var (
		data = struct {
			SceneID          float64 `json:"scene_id"`
			SceneName        string  `json:"scene_name"`
			SceneDescription string  `json:"scene_description"`

			ExecutesIn time.Duration `json:"executes_in"`

			// Audio
			Audio string `json:"audio"`
			// Audio - Clients
			// Audio - Clients - Midi
			AudioMidiEnabled bool   `json:"audio_midi_enabled"`
			AudioMidiClient  uint64 `json:"audio_midi_client"`
			AudioMidiChannel int    `json:"audio_midi_channel"`
			AudioMidiNote    string `json:"audio_midi_note"`
			// Audio - Clients - OSC
			AudioOSCEnabled bool   `json:"audio_osc_enabled"`
			AudioOSCClient  uint64 `json:"audio_osc_client"`
			AudioOSCChannel int    `json:"audio_osc_channel"`
			AudioOSCNote    string `json:"audio_osc_note"`
			// Audio - Clients - GPIO
			AudioGPIOEnabled bool   `json:"audio_gpio_enabled"`
			AudioGPIOClient  uint64 `json:"audio_gpio_client"`
			AudioGPIOChannel int    `json:"audio_gpio_channel"`
			AudioGPIONote    string `json:"audio_gpio_note"`
			// Light
			Light string `json:"light"`
			// Light - Clients
			// Light - Clients - Midi
			LightMidiEnabled bool   `json:"light_midi_enabled"`
			LightMidiClient  uint64 `json:"light_midi_client"`
			LightMidiChannel int    `json:"light_midi_channel"`
			LightMidiNote    string `json:"light_midi_note"`
			// Light - Clients - OSC
			LightOSCEnabled bool   `json:"light_osc_enabled"`
			LightOSCClient  uint64 `json:"light_osc_client"`
			LightOSCChannel int    `json:"light_osc_channel"`
			LightOSCNote    string `json:"light_osc_note"`
			// Light - Clients - GPIO
			LightGPIOEnabled bool   `json:"light_gpio_enabled"`
			LightGPIOClient  uint64 `json:"light_gpio_client"`
			LightGPIOChannel int    `json:"light_gpio_channel"`
			LightGPIONote    string `json:"light_gpio_note"`
			// Video
			Video string `json:"video"`
			// Video - Clients
			// Video - Clients - Midi
			VideoMidiEnabled bool   `json:"video_midi_enabled"`
			VideoMidiClient  uint64 `json:"video_midi_client"`
			VideoMidiChannel int    `json:"video_midi_channel"`
			VideoMidiNote    string `json:"video_midi_note"`
			// Video - Clients - OSC
			VideoOSCEnabled bool   `json:"video_osc_enabled"`
			VideoOSCClient  uint64 `json:"video_osc_client"`
			VideoOSCChannel int    `json:"video_osc_channel"`
			VideoOSCNote    string `json:"video_osc_note"`
			// Video - Clients - GPIO
			VideoGPIOEnabled bool   `json:"video_gpio_enabled"`
			VideoGPIOClient  uint64 `json:"video_gpio_client"`
			VideoGPIOChannel int    `json:"video_gpio_channel"`
			VideoGPIONote    string `json:"video_gpio_note"`
			// Video - Clients - IntraCast
			VideoIntraCastEnabled bool   `json:"video_intra_cast_enabled"`
			VideoIntraCastClient  uint64 `json:"video_intra_cast_client"`
			VideoIntraCastChannel int    `json:"video_intra_cast_channel"`
			VideoIntraCastNote    string `json:"video_intra_cast_note"`
		}{}

		scene database.SceneEntrie
		err   error
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 2, util.Events, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	// Parse&Validate Data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid JSON")
	}
	if data.SceneName == "" && data.SceneID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid Data")
	}

	// General
	scene.SceneID = data.SceneID
	scene.SceneName = data.SceneName
	scene.SceneDescription = data.SceneDescription
	scene.ExecutesIn = data.ExecutesIn

	// Audio
	scene.Audio = data.Audio
	scene.AudioMidiEnabled = data.AudioMidiEnabled
	scene.AudioMidiClient = data.AudioMidiClient
	scene.AudioMidiChannel = data.AudioMidiChannel
	scene.AudioMidiNote = data.AudioMidiNote
	scene.AudioOSCEnabled = data.AudioOSCEnabled
	scene.AudioOSCClient = data.AudioOSCClient
	scene.AudioOSCChannel = data.AudioOSCChannel
	scene.AudioOSCNote = data.AudioOSCNote
	scene.AudioGPIOEnabled = data.AudioGPIOEnabled
	scene.AudioGPIOClient = data.AudioGPIOClient
	scene.AudioGPIOChannel = data.AudioGPIOChannel
	scene.AudioGPIONote = data.AudioGPIONote
	// Light
	scene.Light = data.Light
	scene.LightMidiEnabled = data.LightMidiEnabled
	scene.LightMidiClient = data.LightMidiClient
	scene.LightMidiChannel = data.LightMidiChannel
	scene.LightMidiNote = data.LightMidiNote
	scene.LightOSCEnabled = data.LightOSCEnabled
	scene.LightOSCClient = data.LightOSCClient
	scene.LightOSCChannel = data.LightOSCChannel
	scene.LightOSCNote = data.LightOSCNote
	// Video
	scene.Video = data.Video
	scene.VideoMidiEnabled = data.VideoMidiEnabled
	scene.VideoMidiClient = data.VideoMidiClient
	scene.VideoMidiChannel = data.VideoMidiChannel
	scene.VideoMidiNote = data.VideoMidiNote
	scene.VideoOSCEnabled = data.VideoOSCEnabled
	scene.VideoOSCClient = data.VideoOSCClient
	scene.VideoOSCChannel = data.VideoOSCChannel
	scene.VideoOSCNote = data.VideoOSCNote
	scene.VideoGPIOEnabled = data.VideoGPIOEnabled
	scene.VideoGPIOClient = data.VideoGPIOClient
	scene.VideoGPIOChannel = data.VideoGPIOChannel
	scene.VideoGPIONote = data.VideoGPIONote
	scene.VideoIntraCastEnabled = data.VideoIntraCastEnabled
	scene.VideoIntraCastClient = data.VideoIntraCastClient
	scene.VideoIntraCastChannel = data.VideoIntraCastChannel
	scene.VideoIntraCastNote = data.VideoIntraCastNote

	scene.EventID, _ = strconv.ParseUint(c.Params("id"), 10, 64)

	err = a.DB.Create(&scene).Error
	if err != nil {
		log.Printf("Error creating scene: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	return c.Status(fiber.StatusOK).JSON("")
}

func (a *API) updateScene(c *fiber.Ctx) error {
	var (
		data = struct {
			ID uint `json:"id"`

			SceneID          float64 `json:"scene_id"`
			SceneName        string  `json:"scene_name"`
			SceneDescription string  `json:"scene_description"`

			ExecutesIn time.Duration `json:"executes_in"`

			// Audio
			Audio string `json:"audio"`
			// Audio - Clients
			// Audio - Clients - Midi
			AudioMidiEnabled bool   `json:"audio_midi_enabled"`
			AudioMidiClient  uint64 `json:"audio_midi_client"`
			AudioMidiChannel int    `json:"audio_midi_channel"`
			AudioMidiNote    string `json:"audio_midi_note"`
			// Audio - Clients - OSC
			AudioOSCEnabled bool   `json:"audio_osc_enabled"`
			AudioOSCClient  uint64 `json:"audio_osc_client"`
			AudioOSCChannel int    `json:"audio_osc_channel"`
			AudioOSCNote    string `json:"audio_osc_note"`
			// Audio - Clients - GPIO
			AudioGPIOEnabled bool   `json:"audio_gpio_enabled"`
			AudioGPIOClient  uint64 `json:"audio_gpio_client"`
			AudioGPIOChannel int    `json:"audio_gpio_channel"`
			AudioGPIONote    string `json:"audio_gpio_note"`
			// Light
			Light string `json:"light"`
			// Light - Clients
			// Light - Clients - Midi
			LightMidiEnabled bool   `json:"light_midi_enabled"`
			LightMidiClient  uint64 `json:"light_midi_client"`
			LightMidiChannel int    `json:"light_midi_channel"`
			LightMidiNote    string `json:"light_midi_note"`
			// Light - Clients - OSC
			LightOSCEnabled bool   `json:"light_osc_enabled"`
			LightOSCClient  uint64 `json:"light_osc_client"`
			LightOSCChannel int    `json:"light_osc_channel"`
			LightOSCNote    string `json:"light_osc_note"`
			// Light - Clients - GPIO
			LightGPIOEnabled bool   `json:"light_gpio_enabled"`
			LightGPIOClient  uint64 `json:"light_gpio_client"`
			LightGPIOChannel int    `json:"light_gpio_channel"`
			LightGPIONote    string `json:"light_gpio_note"`
			// Video
			Video string `json:"video"`
			// Video - Clients
			// Video - Clients - Midi
			VideoMidiEnabled bool   `json:"video_midi_enabled"`
			VideoMidiClient  uint64 `json:"video_midi_client"`
			VideoMidiChannel int    `json:"video_midi_channel"`
			VideoMidiNote    string `json:"video_midi_note"`
			// Video - Clients - OSC
			VideoOSCEnabled bool   `json:"video_osc_enabled"`
			VideoOSCClient  uint64 `json:"video_osc_client"`
			VideoOSCChannel int    `json:"video_osc_channel"`
			VideoOSCNote    string `json:"video_osc_note"`
			// Video - Clients - GPIO
			VideoGPIOEnabled bool   `json:"video_gpio_enabled"`
			VideoGPIOClient  uint64 `json:"video_gpio_client"`
			VideoGPIOChannel int    `json:"video_gpio_channel"`
			VideoGPIONote    string `json:"video_gpio_note"`
			// Video - Clients - IntraCast
			VideoIntraCastEnabled bool   `json:"video_intra_cast_enabled"`
			VideoIntraCastClient  uint64 `json:"video_intra_cast_client"`
			VideoIntraCastChannel int    `json:"video_intra_cast_channel"`
			VideoIntraCastNote    string `json:"video_intra_cast_note"`
		}{}

		scene database.SceneEntrie

		err error
	)
	if !util.CheckPermissions(c.GetReqHeaders(), 2, util.Events, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}
	// Parse body
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Error parsing data: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	if data.ID == 0 || data.SceneID == 0 || data.SceneName == "" {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid data")
	}

	// Find Scene
	err = a.DB.Find(&scene, data.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Error finding scene: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	// Update the values
	// General
	scene.SceneID = data.SceneID
	scene.SceneName = data.SceneName
	scene.SceneDescription = data.SceneDescription
	scene.ExecutesIn = data.ExecutesIn

	// Audio
	scene.Audio = data.Audio
	scene.AudioMidiEnabled = data.AudioMidiEnabled
	scene.AudioMidiClient = data.AudioMidiClient
	scene.AudioMidiChannel = data.AudioMidiChannel
	scene.AudioMidiNote = data.AudioMidiNote
	scene.AudioOSCEnabled = data.AudioOSCEnabled
	scene.AudioOSCClient = data.AudioOSCClient
	scene.AudioOSCChannel = data.AudioOSCChannel
	scene.AudioOSCNote = data.AudioOSCNote
	scene.AudioGPIOEnabled = data.AudioGPIOEnabled
	scene.AudioGPIOClient = data.AudioGPIOClient
	scene.AudioGPIOChannel = data.AudioGPIOChannel
	scene.AudioGPIONote = data.AudioGPIONote
	// Light
	scene.Light = data.Light
	scene.LightMidiEnabled = data.LightMidiEnabled
	scene.LightMidiClient = data.LightMidiClient
	scene.LightMidiChannel = data.LightMidiChannel
	scene.LightMidiNote = data.LightMidiNote
	scene.LightOSCEnabled = data.LightOSCEnabled
	scene.LightOSCClient = data.LightOSCClient
	scene.LightOSCChannel = data.LightOSCChannel
	scene.LightOSCNote = data.LightOSCNote
	// Video
	scene.Video = data.Video
	scene.VideoMidiEnabled = data.VideoMidiEnabled
	scene.VideoMidiClient = data.VideoMidiClient
	scene.VideoMidiChannel = data.VideoMidiChannel
	scene.VideoMidiNote = data.VideoMidiNote
	scene.VideoOSCEnabled = data.VideoOSCEnabled
	scene.VideoOSCClient = data.VideoOSCClient
	scene.VideoOSCChannel = data.VideoOSCChannel
	scene.VideoOSCNote = data.VideoOSCNote
	scene.VideoGPIOEnabled = data.VideoGPIOEnabled
	scene.VideoGPIOClient = data.VideoGPIOClient
	scene.VideoGPIOChannel = data.VideoGPIOChannel
	scene.VideoGPIONote = data.VideoGPIONote
	scene.VideoIntraCastEnabled = data.VideoIntraCastEnabled
	scene.VideoIntraCastClient = data.VideoIntraCastClient
	scene.VideoIntraCastChannel = data.VideoIntraCastChannel
	scene.VideoIntraCastNote = data.VideoIntraCastNote

	err = a.DB.Save(scene).Error
	if err != nil {
		log.Printf("Error saving scene: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	return c.Status(fiber.StatusOK).JSON("Successfully updated scene")
}

func (a *API) deleteScene(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 3, util.Events, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	err := a.DB.Delete(&database.SceneEntrie{}, c.Params("id")).Error
	if err != nil {
		log.Printf("Error deleting scene: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	return c.Status(fiber.StatusOK).JSON("")
}

func (a *API) testScene(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 2, util.Events, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	var scene database.SceneEntrie
	err := a.DB.First(&scene, c.Params("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Error getting scene: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	log.Printf("Testing scene (%s): %d\n", c.Params("id"), scene.ID)

	return c.Status(fiber.StatusOK).JSON("W.I.P.")
}

func (a *API) updateActiveScene(c *fiber.Ctx) error {
	if !util.CheckPermissions(c.GetReqHeaders(), 2, util.Events, a.DB) {
		return c.Status(fiber.StatusForbidden).JSON("")
	}

	var event database.Event
	err := a.DB.First(&event, c.Params("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("")
		}
		log.Printf("Error getting event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	event.CurrentScene, _ = c.ParamsInt("sceneID")

	err = a.DB.Save(&event).Error
	if err != nil {
		log.Printf("Error updating event: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON("")
	}

	return c.Status(fiber.StatusOK).JSON("Selected Scene updated successfully")
}
