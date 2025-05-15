package v1_test

import (
	"net/http"
	"testing"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestCreateApp(t *testing.T) {
	// Test case 1: Create app with invalid input (missing required fields)
	t.Run("create app with invalid input", func(t *testing.T) {
		req := structs.App{
			App: "MyTestApp", // missing other required fields
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/apps")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	// Test case 2: Create app with valid input
	t.Run("create app with valid input", func(t *testing.T) {
		req := structs.App{
			App:           "MyTestApp",
			Category:      "Utilities",
			Rating:        4.5,
			Reviews:       1000,
			Size:          "15MB",
			Installs:      "50000",
			Type:          "Free",
			Price:         "$0",
			ContentRating: "Everyone",
			Genres:        "Tools",
			LastUpdated:   "2025-05-14",
			CurrentVer:    "1.0.0",
			AndroidVer:    "5.0 and up",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/apps")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())
	})

	// Test case 3: Create app with missing required fields (e.g., app name)
	t.Run("create app with missing fields", func(t *testing.T) {
		req := structs.App{
			Category:      "Utilities",
			Rating:        4.5,
			Reviews:       1000,
			Size:          "15MB",
			Installs:      "50000",
			Type:          "Free",
			Price:         "$0",
			ContentRating: "Everyone",
			Genres:        "Tools",
			LastUpdated:   "2025-05-14",
			CurrentVer:    "1.0.0",
			AndroidVer:    "5.0 and up",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/apps")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	// Cleanup after test
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM apps WHERE app='MyTestApp'")
		assert.Nil(t, err)
	})
}
func TestGetApps(t *testing.T) {
	// Test case 1: Get apps with valid limit and offset
	t.Run("get apps with valid limit and offset", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/apps?limit=10&offset=0")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	// Test case 2: Get apps with invalid (negative) limit
	t.Run("get apps with invalid limit", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/apps?limit=-10&offset=0")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	// Test case 3: Get apps with invalid offset (negative offset)
	t.Run("get apps with invalid offset", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/apps?limit=10&offset=-5")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
}
func TestGetAppById(t *testing.T) {
	// Test case 1: Get app by valid ID
	t.Run("get app by valid  ID", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/apps/2")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	// Test case 2: Get app by non-existing ID
	t.Run("get app by non-existing ID", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/apps/99999") // Non-existent ID

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}
func TestDeleteApp(t *testing.T) {
	// Test case 1: Delete app by valid ID
	t.Run("delete app by valid ID", func(t *testing.T) {
		// Assume the app with ID 1 exists
		res, err := client.
			R().
			EnableTrace().
			Delete("/api/v1/apps/2")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	// Test case 2: Delete app by non-existing ID
	t.Run("delete app by non-existing ID", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Delete("/api/v1/apps/99999") // Non-existent ID

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}
func TestUpdateApp(t *testing.T) {
	// Test case 1: Update app with valid ID and valid data
	t.Run("update app with valid data", func(t *testing.T) {
		req := structs.App{
			App:           "UpdatedTestApp",
			Category:      "Games",
			Rating:        4.7,
			Reviews:       1200,
			Size:          "20MB",
			Installs:      "60000",
			Type:          "Paid",
			Price:         "$2.99",
			ContentRating: "Teen",
			Genres:        "Action",
			LastUpdated:   "2025-06-14",
			CurrentVer:    "1.1.0",
			AndroidVer:    "6.0 and up",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/apps/1") // Assume app ID 1 exists

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	// Test case 2: Update app with missing required field (e.g., app name)
	t.Run("update app with missing required fields", func(t *testing.T) {
		req := structs.App{
			Category:      "Games",
			Rating:        4.7,
			Reviews:       1200,
			Size:          "20MB",
			Installs:      "60000",
			Type:          "Paid",
			Price:         "$2.99",
			ContentRating: "Teen",
			Genres:        "Action",
			LastUpdated:   "2025-06-14",
			CurrentVer:    "1.1.0",
			AndroidVer:    "6.0 and up",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/apps/3")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	// Test case 3: Update app with non-existing ID
	t.Run("update app with non-existing ID", func(t *testing.T) {
		req := structs.App{
			App:           "UpdatedTestApp",
			Category:      "Games",
			Rating:        4.7,
			Reviews:       1200,
			Size:          "20MB",
			Installs:      "60000",
			Type:          "Paid",
			Price:         "$2.99",
			ContentRating: "Teen",
			Genres:        "Action",
			LastUpdated:   "2025-06-14",
			CurrentVer:    "1.1.0",
			AndroidVer:    "6.0 and up",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/apps/99999") // Non-existent ID

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}
