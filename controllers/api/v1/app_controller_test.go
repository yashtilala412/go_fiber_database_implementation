package v1_test

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	// Assuming your models and utils packages are importable
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// client and db are global variables provided by TestMain in v1_test.go

// Helper function to create a default valid App struct for requests
func newTestApp(tag string) models.App {
	return models.App{
		App:           "Test App " + tag,
		Category:      "TEST_CATEGORY",
		Rating:        4.5,
		Reviews:       100,
		Size:          "10M",
		Installs:      "1,000+",
		Type:          "Free",
		Price:         "0",
		ContentRating: "Everyone",
		Genres:        "Testing",
		LastUpdated:   time.Now().Format("Jan 02, 2006"),
		CurrentVer:    "1.0",
		AndroidVer:    "4.0.3 and up",
	}
}

// TestAppController_CreateApp tests the POST /api/v1/apps endpoint
func TestAppController_CreateApp(t *testing.T) {
	t.Run("create app with valid input", func(t *testing.T) {
		reqBody := newTestApp(time.Now().Format("20060102150405-valid"))
		var resBody utils.JSONResponse
		var createdAppID int

		// Cleanup: Delete the app created in this test
		t.Cleanup(func() {
			if createdAppID != 0 {
				log.Printf("Cleaning up app created in TestAppController_CreateApp (valid): %d", createdAppID)
				client.R().Delete(fmt.Sprintf("/api/v1/apps/%d", createdAppID)) // Best effort cleanup
			}
		})

		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Post("/api/v1/apps")

		require.Nil(t, err, "Error making POST request")
		assert.Equal(t, http.StatusCreated, res.StatusCode(), "Expected 201 Created status code")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success'")

		// Access data directly from resBody after SetResult
		responseDataMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map[string]interface{}")

		idFloat, ok := responseDataMap["id"].(float64)
		require.True(t, ok, "Could not find 'id' field in response data or it's not a number")
		createdAppID = int(idFloat) // Store ID for cleanup

		assert.NotEqual(t, 0, createdAppID, "Created App ID should not be zero")
		assert.Equal(t, reqBody.App, responseDataMap["app"], "App name in response does not match")
	})

	t.Run("create app with invalid input (missing required field)", func(t *testing.T) {
		reqBody := models.App{
			Price: "100",
		}
		var resBody utils.JSONResponse

		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Post("/api/v1/apps")

		require.Nil(t, err, "Error making POST request with invalid input")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})
}

// TestAppController_GetApp tests the GET /api/v1/apps/{id} endpoint
func TestAppController_GetApp(t *testing.T) {
	var createdAppID int
	reqBody := newTestApp(time.Now().Format("20060102150405-getbyid"))
	var setupResBody utils.JSONResponse // Declare variable for setup response body

	// Setup: Create an app to fetch
	setupRes, setupErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&setupResBody). // Resty attempts to unmarshal here
		Post("/api/v1/apps")
	require.Nil(t, setupErr)
	require.Equal(t, http.StatusCreated, setupRes.StatusCode)

	// Access data directly from setupResBody after SetResult
	setupDataMap, ok := setupResBody.Data.(map[string]interface{})
	require.True(t, ok)
	idFloat, ok := setupDataMap["id"].(float64)
	require.True(t, ok)
	createdAppID = int(idFloat)
	require.NotEqual(t, 0, createdAppID)

	// Cleanup: Delete the app created for this test
	t.Cleanup(func() {
		log.Printf("Cleaning up app created for TestAppController_GetApp: %d", createdAppID)
		client.R().Delete(fmt.Sprintf("/api/v1/apps/%d", createdAppID)) // Best effort cleanup
	})

	t.Run("get app with valid id", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get(fmt.Sprintf("/api/v1/apps/%d", createdAppID))

		require.Nil(t, err, "Error making GET request")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success'")

		retrievedAppMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map[string]interface{}")
		assert.Equal(t, float64(createdAppID), retrievedAppMap["id"], "Retrieved app ID does not match")
		assert.Equal(t, reqBody.App, retrievedAppMap["app"], "Retrieved app name does not match")
	})

	t.Run("get app with non-existent id (expect 404)", func(t *testing.T) {
		res, err := client.R().Get("/api/v1/apps/99999") // Use a non-existent ID

		require.Nil(t, err, "Error making GET request for non-existent id")
		assert.Equal(t, http.StatusNotFound, res.StatusCode(), "Expected 404 Not Found status code")
	})

	t.Run("get app with invalid id format (expect 400)", func(t *testing.T) {
		res, err := client.R().Get("/api/v1/apps/abc") // Use invalid ID format

		require.Nil(t, err, "Error making GET request with invalid id format")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
	})
}

// TestAppController_GetApps tests the GET /api/v1/apps endpoint including pagination
func TestAppController_GetApps(t *testing.T) {
	numInitialApps := 8 // Number of apps to create for pagination testing
	var initialAppIDs []int

	// Setup: Create initial apps for pagination tests
	log.Printf("Creating %d initial apps for TestAppController_GetApps...", numInitialApps)
	for i := 0; i < numInitialApps; i++ {
		reqBody := newTestApp(fmt.Sprintf("list-%s-%d", time.Now().Format("20060102150405"), i))
		var resBody utils.JSONResponse // Declare variable for setup response body
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Post("/api/v1/apps")

		require.Nil(t, err, fmt.Sprintf("Failed to create initial app %d for list test", i))
		require.Equal(t, http.StatusCreated, res.StatusCode(), fmt.Sprintf("Expected 201 Created for initial app %d list test", i))

		// Access data directly from resBody after SetResult
		responseDataMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, fmt.Sprintf("Initial app %d list test response data is not a map", i))
		idFloat, ok := responseDataMap["id"].(float64)
		require.True(t, ok, fmt.Sprintf("Could not find 'id' for initial app %d list test or it's not a number", i))
		initialAppIDs = append(initialAppIDs, int(idFloat))
	}
	log.Printf("Finished creating %d initial apps for TestAppController_GetApps. IDs: %v", numInitialApps, initialAppIDs)

	// Cleanup: Delete the initial apps created for this test
	t.Cleanup(func() {
		log.Printf("Cleaning up %d initial apps from TestAppController_GetApps...", len(initialAppIDs))
		for _, appID := range initialAppIDs {
			if appID != 0 {
				client.R().Delete(fmt.Sprintf("/api/v1/apps/%d", appID)) // Best effort cleanup
			}
		}
	})

	t.Run("list apps (basic)", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request (basic list)")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code (basic list)")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' (basic list)")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for basic list")
		assert.Len(t, resBody.Data, numInitialApps, fmt.Sprintf("Expected %d items in the basic list", numInitialApps))
	})

	t.Run("list apps with limit", func(t *testing.T) {
		limit := 3
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("limit", fmt.Sprintf("%d", limit)).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request with limit")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code with limit")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' with limit")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for list with limit")
		assert.Len(t, resBody.Data, limit, fmt.Sprintf("Expected %d items with limit %d", limit, limit))
	})

	t.Run("list apps with offset", func(t *testing.T) {
		offset := 2
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("offset", fmt.Sprintf("%d", offset)).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request with offset")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code with offset")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' with offset")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for list with offset")
		expectedCount := numInitialApps - offset
		if expectedCount < 0 {
			expectedCount = 0
		}
		assert.Len(t, resBody.Data, expectedCount, fmt.Sprintf("Expected %d items with offset %d", expectedCount, offset))
	})

	t.Run("list apps with limit and offset", func(t *testing.T) {
		limit := 3
		offset := 4
		var resBody utils.JSONResponse
		res, err := client.R().
			SetQueryParams(map[string]string{
				"limit": fmt.Sprintf("%d", limit), "offset": fmt.Sprintf("%d", offset),
			}).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request with limit and offset")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code with limit and offset")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' with limit and offset")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for list with limit and offset")

		remainingAfterOffset := numInitialApps - offset
		expectedCount := limit
		if remainingAfterOffset < limit {
			expectedCount = remainingAfterOffset
		}
		if expectedCount < 0 {
			expectedCount = 0
		}
		assert.Len(t, resBody.Data, expectedCount, fmt.Sprintf("Expected %d items with limit %d and offset %d", expectedCount, limit, offset))
	})

	t.Run("list apps with invalid limit (expect 400)", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("limit", "invalid").
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request with invalid limit")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request with invalid limit")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error' with invalid limit")
	})

	t.Run("list apps with invalid offset (expect 400)", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("offset", "invalid").
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request with invalid offset")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request with invalid offset")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error' with invalid offset")
	})

	t.Run("list apps with limit exceeding MaxLimit (expect 400)", func(t *testing.T) {
		maxLimit := 500
		limitExceeding := maxLimit + 1
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("limit", fmt.Sprintf("%d", limitExceeding)).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Get("/api/v1/apps")

		require.Nil(t, err, "Error making GET request with limit exceeding MaxLimit")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request with limit exceeding MaxLimit")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error' with limit exceeding MaxLimit")
	})
}

// TestAppController_UpdateApp tests the PUT /api/v1/apps/{id} endpoint
func TestAppController_UpdateApp(t *testing.T) {
	var createdAppID int
	reqBody := newTestApp(time.Now().Format("20060102150405-update"))
	var setupResBody utils.JSONResponse // Declare variable for setup response body

	// Setup: Create an app to update
	setupRes, setupErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&setupResBody). // Resty attempts to unmarshal here
		Post("/api/v1/apps")
	require.Nil(t, setupErr)
	require.Equal(t, http.StatusCreated, setupRes.StatusCode)

	// Access data directly from setupResBody after SetResult
	setupDataMap, ok := setupResBody.Data.(map[string]interface{})
	require.True(t, ok)
	idFloat, ok := setupDataMap["id"].(float64)
	require.True(t, ok)
	createdAppID = int(idFloat)
	require.NotEqual(t, 0, createdAppID)

	// Cleanup: Delete the app created for this test
	t.Cleanup(func() {
		log.Printf("Cleaning up app created for TestAppController_UpdateApp: %d", createdAppID)
		client.R().Delete(fmt.Sprintf("/api/v1/apps/%d", createdAppID)) // Best effort cleanup
	})

	t.Run("update app with valid input", func(t *testing.T) {
		updatedAppName := reqBody.App + " Updated"
		updatedReqBody := newTestApp(time.Now().Format("20060102150405-updated"))
		updatedReqBody.App = updatedAppName // Modify the name
		updatedReqBody.Rating = 4.9         // Modify another field

		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(updatedReqBody).
			SetResult(&resBody). // Resty attempts to unmarshal here
			Put(fmt.Sprintf("/api/v1/apps/%d", createdAppID))

		require.Nil(t, err, "Error making PUT request")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success'")

		updatedAppMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map[string]interface{}")
		assert.Equal(t, updatedAppName, updatedAppMap["app"], "Updated app name in response does not match")
		assert.Equal(t, 4.9, updatedAppMap["rating"], "Updated rating in response does not match")
	})

	t.Run("update app with non-existent id (expect 404)", func(t *testing.T) {
		updatedReqBody := newTestApp(time.Now().Format("20060102150405-nonexistent-update"))
		res, err := client.R().SetHeader("Content-Type", "application/json").SetBody(updatedReqBody).Put("/api/v1/apps/99999")

		require.Nil(t, err, "Error making PUT request for non-existent id")
		assert.Equal(t, http.StatusNotFound, res.StatusCode(), "Expected 404 Not Found status code")
	})

	t.Run("update app with invalid input (missing required field - expect 400)", func(t *testing.T) {
		invalidReqBody := models.App{
			Price: "100",
		}
		var resBody utils.JSONResponse
		res, err := client.R().SetHeader("Content-Type", "application/json").SetBody(invalidReqBody).SetResult(&resBody).Put(fmt.Sprintf("/api/v1/apps/%d", createdAppID))

		require.Nil(t, err, "Error making PUT request with invalid input")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})
	t.Run("update app with invalid id format (expect 400)", func(t *testing.T) {
		updatedReqBody := newTestApp(time.Now().Format("20060102150405-invalidformat-update"))
		res, err := client.R().SetHeader("Content-Type", "application/json").SetBody(updatedReqBody).Put("/api/v1/apps/abc")

		require.Nil(t, err, "Error making PUT request with invalid id format")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
	})
}

// TestAppController_DeleteApp tests the DELETE /api/v1/apps/{id} endpoint
func TestAppController_DeleteApp(t *testing.T) {
	var createdAppID int
	reqBody := newTestApp(time.Now().Format("20060102150405-delete"))
	var setupResBody utils.JSONResponse // Declare variable for setup response body

	// Setup: Create an app to delete
	setupRes, setupErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&setupResBody). // Resty attempts to unmarshal here
		Post("/api/v1/apps")
	require.Nil(t, setupErr)
	require.Equal(t, http.StatusCreated, setupRes.StatusCode)

	// Access data directly from setupResBody after SetResult
	setupDataMap, ok := setupResBody.Data.(map[string]interface{})
	require.True(t, ok)
	idFloat, ok := setupDataMap["id"].(float64)
	require.True(t, ok)
	createdAppID = int(idFloat)
	require.NotEqual(t, 0, createdAppID)

	t.Run("delete app with valid id", func(t *testing.T) {
		res, err := client.R().Delete(fmt.Sprintf("/api/v1/apps/%d", createdAppID))

		require.Nil(t, err, "Error making DELETE request")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code")
		// Optional: Assert on the success response body if it contains a message
		// The AppController returns constants.AppsDeletedSuccessfully on success.
		// You might want to assert the response body contains this string or a specific JSON structure.
	})

	t.Run("delete app with non-existent id (expect 404)", func(t *testing.T) {
		res, err := client.R().Delete("/api/v1/apps/99999") // Use a non-existent ID

		require.Nil(t, err, "Error making DELETE request for non-existent id")
		assert.Equal(t, http.StatusNotFound, res.StatusCode(), "Expected 404 Not Found status code")
		// Optional: Assert on the error response body structure if needed
	})
	t.Run("delete app with invalid id format (expect 400)", func(t *testing.T) {
		res, err := client.R().Delete("/api/v1/apps/abc")

		require.Nil(t, err, "Error making DELETE request with invalid id format")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
		// Optional: Assert on the error response body structure if needed
	})
}
