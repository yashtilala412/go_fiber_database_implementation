package v1_test

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	// Assuming your models and utils packages are importable
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models" // Import your actual models
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// client and db are global variables provided by TestMain in v1_test.go

// Helper to create a default valid Review struct for requests
func newTestReview(tag string) models.Review {
	// Use fields matching your models.Review struct and validation tags
	return models.Review{
		App:                   "Test App Name " + tag, // Use App (string) field from your model
		TranslatedReview:      "This is a test review. " + tag,
		Sentiment:             "Positive",
		SentimentPolarity:     models.NullableFloat64{Float64: 0.5, Valid: true}, // Use NullableFloat64 as defined in your models
		SentimentSubjectivity: models.NullableFloat64{Float64: 0.8, Valid: true}, // Use NullableFloat64 as defined in your models
		// Rating field is NOT included as it's not in your models.Review struct
	}
}

// TestReviewController_CreateReview tests the POST /api/v1/reviews endpoint
func TestReviewController_CreateReview(t *testing.T) {
	t.Run("create review with valid input", func(t *testing.T) {
		reqBody := newTestReview("valid")
		var resBody utils.JSONResponse
		var createdReviewID int

		// Cleanup: Delete the review created in this test
		t.Cleanup(func() {
			if createdReviewID != 0 {
				log.Printf("Cleaning up review created in TestReviewController_CreateReview (valid): %d", createdReviewID)
				client.R().Delete(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID)) // Best effort cleanup
			}
		})

		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&resBody).
			Post("/api/v1/reviews")

		require.Nil(t, err, "Error making POST request to /api/v1/reviews")
		assert.Equal(t, http.StatusCreated, res.StatusCode(), "Expected 201 Created status code")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success'")

		responseDataMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map[string]interface{}")

		idFloat, ok := responseDataMap["id"].(float64)
		require.True(t, ok, "Could not find 'id' field in response data or it's not a number")
		createdReviewID = int(idFloat) // Store ID for cleanup

		assert.NotEqual(t, 0, createdReviewID, "Created Review ID should not be zero")
		assert.Equal(t, reqBody.App, responseDataMap["app"], "App name in response does not match")
		assert.Equal(t, reqBody.TranslatedReview, responseDataMap["translated_review"], "Review text in response does not match")

		// Checking NullableFloat64 fields
		// Check SentimentPolarity
		valuePolarity, okPolarity := responseDataMap["sentiment_polarity"]
		if reqBody.SentimentPolarity.Valid {
			// Declare and use sentimentPolarity ONLY in this block
			sentimentPolarity, isFloatPolarity := valuePolarity.(float64)
			require.True(t, okPolarity, "Expected 'sentiment_polarity' field in response")
			require.True(t, isFloatPolarity, "Expected 'sentiment_polarity' value to be a float64")
			assert.Equal(t, reqBody.SentimentPolarity.Float64, sentimentPolarity, "Sentiment polarity in response does not match")
		} else {
			_, existsPolarity := responseDataMap["sentiment_polarity"]
			assert.False(t, existsPolarity, "Did not expect 'sentiment_polarity' field in response if not valid in request")
		}

		// Check SentimentSubjectivity
		valueSubjectivity, okSubjectivity := responseDataMap["sentiment_subjectivity"]
		if reqBody.SentimentSubjectivity.Valid {
			// Declare and use sentimentSubjectivity ONLY in this block
			sentimentSubjectivity, isFloatSubjectivity := valueSubjectivity.(float64)
			require.True(t, okSubjectivity, "Expected 'sentiment_subjectivity' field in response")
			require.True(t, isFloatSubjectivity, "Expected 'sentiment_subjectivity' value to be a float64")
			assert.Equal(t, reqBody.SentimentSubjectivity.Float64, sentimentSubjectivity, "Sentiment subjectivity in response does not match")
		} else {
			_, existsSubjectivity := responseDataMap["sentiment_subjectivity"]
			assert.False(t, existsSubjectivity, "Did not expect 'sentiment_subjectivity' field in response if not valid in request")
		}
	})

	t.Run("create review with invalid input (missing required field)", func(t *testing.T) {
		reqBody := models.Review{
			// Missing required fields like App, TranslatedReview, Sentiment
			SentimentPolarity: models.NullableFloat64{Float64: 0.1, Valid: true},
		}
		var resBody utils.JSONResponse

		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&resBody).
			Post("/api/v1/reviews")

		require.Nil(t, err, "Error making POST request with invalid input")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})
}

// TestReviewController_GetReview tests the GET /api/v1/reviews/{id} endpoint
func TestReviewController_GetReview(t *testing.T) {
	var createdReviewID int
	reqBody := newTestReview("getbyid")
	var setupResBody utils.JSONResponse

	// Setup: Create a review to fetch
	setupRes, setupErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&setupResBody).
		Post("/api/v1/reviews")
	require.Nil(t, setupErr)
	require.Equal(t, http.StatusCreated, setupRes.StatusCode, "Setup: Expected 201 Created when creating review for GetReview test")

	setupDataMap, ok := setupResBody.Data.(map[string]interface{})
	require.True(t, ok)
	idFloat, ok := setupDataMap["id"].(float64)
	require.True(t, ok, "Could not find 'id' field in setup response data or it's not a number")
	createdReviewID = int(idFloat)
	require.NotEqual(t, 0, createdReviewID)

	// Cleanup: Delete the review created for this test
	t.Cleanup(func() {
		log.Printf("Cleaning up review created for TestReviewController_GetReview: %d", createdReviewID)
		client.R().Delete(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID)) // Best effort cleanup
	})

	t.Run("get review with valid id", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetResult(&resBody).
			Get(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID))

		require.Nil(t, err, "Error making GET request")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success'")

		retrievedReviewMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map[string]interface{}")
		assert.Equal(t, float64(createdReviewID), retrievedReviewMap["id"], "Retrieved review ID does not match")
		assert.Equal(t, reqBody.App, retrievedReviewMap["app"], "Retrieved review app name does not match")
		assert.Equal(t, reqBody.TranslatedReview, retrievedReviewMap["translated_review"], "Retrieved review text does not match")
		// Optional: Check NullableFloat64 fields - can add similar checks as in Create test if needed
	})

	t.Run("get review with non-existent id (expect 404)", func(t *testing.T) {
		res, err := client.R().Get("/api/v1/reviews/99999") // Use a non-existent ID

		require.Nil(t, err, "Error making GET request for non-existent id")
		assert.Equal(t, http.StatusNotFound, res.StatusCode(), "Expected 404 Not Found status code")
	})

	t.Run("get review with invalid id format (expect 400)", func(t *testing.T) {
		res, err := client.R().Get("/api/v1/reviews/abc") // Use invalid ID format

		require.Nil(t, err, "Error making GET request with invalid id format")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
	})
}

// TestReviewController_GetReviews tests the GET /api/v1/reviews endpoint including pagination
func TestReviewController_GetReviews(t *testing.T) {
	numInitialReviews := 8
	var initialReviewIDs []int

	// Setup: Create initial reviews for pagination tests
	log.Printf("Creating %d initial reviews for TestReviewController_GetReviews...", numInitialReviews)
	for i := 0; i < numInitialReviews; i++ {
		reqBody := newTestReview(fmt.Sprintf("list-%s-%d", time.Now().Format("20060102150405"), i))
		var resBody utils.JSONResponse
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(reqBody).
			SetResult(&resBody).
			Post("/api/v1/reviews")

		require.Nil(t, err, fmt.Sprintf("Failed to create initial review %d for list test", i))
		require.Equal(t, http.StatusCreated, res.StatusCode(), fmt.Sprintf("Setup: Expected 201 Created for initial review %d list test", i))

		responseDataMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, fmt.Sprintf("Initial review %d list test response data is not a map", i))
		idFloat, ok := responseDataMap["id"].(float64)
		require.True(t, ok, fmt.Sprintf("Could not find 'id' for initial review %d list test or it's not a number", i))
		initialReviewIDs = append(initialReviewIDs, int(idFloat))
	}
	log.Printf("Finished creating %d initial reviews for TestReviewController_GetReviews. IDs: %v", numInitialReviews, initialReviewIDs)

	// Cleanup: Delete the initial reviews created for this test
	t.Cleanup(func() {
		log.Printf("Cleaning up %d initial reviews from TestReviewController_GetReviews...", len(initialReviewIDs))
		for _, reviewID := range initialReviewIDs {
			if reviewID != 0 {
				client.R().Delete(fmt.Sprintf("/api/v1/reviews/%d", reviewID))
			}
		}
	})

	t.Run("list reviews (basic)", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request (basic list)")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code (basic list)")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' (basic list)")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for basic list")
		assert.Len(t, resBody.Data, numInitialReviews, fmt.Sprintf("Expected %d items in the basic list", numInitialReviews))
	})

	t.Run("list reviews with limit", func(t *testing.T) {
		limit := 3
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("limit", fmt.Sprintf("%d", limit)).
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request with limit")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code with limit")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' with limit")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for list with limit")
		assert.Len(t, resBody.Data, limit, fmt.Sprintf("Expected %d items with limit %d", limit, limit))
	})

	t.Run("list reviews with offset", func(t *testing.T) {
		offset := 2
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("offset", fmt.Sprintf("%d", offset)).
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request with offset")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code with offset")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' with offset")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for list with offset")
		expectedCount := numInitialReviews - offset
		if expectedCount < 0 {
			expectedCount = 0
		}
		assert.Len(t, resBody.Data, expectedCount, fmt.Sprintf("Expected %d items with offset %d", expectedCount, offset))
	})

	t.Run("list reviews with limit and offset", func(t *testing.T) {
		limit := 3
		offset := 4
		var resBody utils.JSONResponse
		res, err := client.R().
			SetQueryParams(map[string]string{
				"limit": fmt.Sprintf("%d", limit), "offset": fmt.Sprintf("%d", offset),
			}).
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request with limit and offset")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code with limit and offset")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success' with limit and offset")

		responseDataValue := reflect.ValueOf(resBody.Data)
		require.Containsf(t, []reflect.Kind{reflect.Slice, reflect.Array}, responseDataValue.Kind(), "Response data is not a slice/array for list with limit and offset")

		remainingAfterOffset := numInitialReviews - offset
		expectedCount := limit
		if remainingAfterOffset < limit {
			expectedCount = remainingAfterOffset
		}
		if expectedCount < 0 {
			expectedCount = 0
		}
		assert.Len(t, resBody.Data, expectedCount, fmt.Sprintf("Expected %d items with limit %d and offset %d", expectedCount, limit, offset))
	})

	t.Run("list reviews with invalid limit (expect 400)", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("limit", "invalid").
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request with invalid limit")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request with invalid limit")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})

	t.Run("list reviews with invalid offset (expect 400)", func(t *testing.T) {
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("offset", "invalid").
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request with invalid offset")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request with invalid offset")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})

	t.Run("list reviews with limit exceeding MaxLimit (expect 400)", func(t *testing.T) {
		maxLimit := 500
		limitExceeding := maxLimit + 1
		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetQueryParam("limit", fmt.Sprintf("%d", limitExceeding)).
			SetResult(&resBody).
			Get("/api/v1/reviews")

		require.Nil(t, err, "Error making GET request with limit exceeding MaxLimit")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request with limit exceeding MaxLimit")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})
}

// TestReviewController_UpdateReview tests the PUT /api/v1/reviews/{id} endpoint
func TestReviewController_UpdateReview(t *testing.T) {
	var createdReviewID int
	reqBody := newTestReview("update")
	var setupResBody utils.JSONResponse

	// Setup: Create a review to update
	setupRes, setupErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&setupResBody).
		Post("/api/v1/reviews")
	require.Nil(t, setupErr)
	require.Equal(t, http.StatusCreated, setupRes.StatusCode, "Setup: Expected 201 Created when creating review for UpdateReview test")

	setupDataMap, ok := setupResBody.Data.(map[string]interface{})
	require.True(t, ok)
	idFloat, ok := setupDataMap["id"].(float64)
	require.True(t, ok)
	createdReviewID = int(idFloat)
	require.NotEqual(t, 0, createdReviewID)

	// Cleanup: Delete the review created for this test
	t.Cleanup(func() {
		log.Printf("Cleaning up review created for TestReviewController_UpdateReview: %d", createdReviewID)
		client.R().Delete(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID)) // Best effort cleanup
	})

	t.Run("update review with valid input", func(t *testing.T) {
		updatedReviewText := reqBody.TranslatedReview + " Updated"
		updatedReqBody := newTestReview("updated") // Create a new review body
		// Do NOT set ReviewId in the request body unless your API explicitly requires it and uses it.
		// The update target is specified by the path parameter.
		updatedReqBody.TranslatedReview = updatedReviewText // Modify the text
		updatedReqBody.App = "Updated App Name"             // Modify the App string
		// Optional: Modify NullableFloat64 fields
		updatedReqBody.SentimentPolarity = models.NullableFloat64{Float64: -0.1, Valid: true}

		var resBody utils.JSONResponse
		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(updatedReqBody).
			SetResult(&resBody).
			Put(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID)) // Use the ID from the path

		require.Nil(t, err, "Error making PUT request")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code")
		assert.Equal(t, "success", resBody.Status, "Expected JSON response status to be 'success'")

		updatedReviewMap, ok := resBody.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map[string]interface{}")
		assert.Equal(t, float64(createdReviewID), updatedReviewMap["id"], "Updated review ID in response does not match")
		assert.Equal(t, updatedReqBody.App, updatedReviewMap["app"], "Updated review app name in response does not match")
		assert.Equal(t, updatedReviewText, updatedReviewMap["translated_review"], "Updated review text in response does not match")
		// Checking updated NullableFloat64 fields - declared inside if block
		valuePolarity, okPolarity := updatedReviewMap["sentiment_polarity"]
		if updatedReqBody.SentimentPolarity.Valid {
			sentimentPolarity, isFloatPolarity := valuePolarity.(float64)
			require.True(t, okPolarity, "Expected 'sentiment_polarity' field in updated response")
			require.True(t, isFloatPolarity, "Expected 'sentiment_polarity' value to be a float64 in updated response")
			assert.Equal(t, updatedReqBody.SentimentPolarity.Float64, sentimentPolarity, "Updated sentiment polarity in response does not match")
		} else {
			_, existsPolarity := updatedReviewMap["sentiment_polarity"]
			assert.False(t, existsPolarity, "Did not expect 'sentiment_polarity' field in updated response if not valid in request")
		}

		valueSubjectivity, okSubjectivity := updatedReviewMap["sentiment_subjectivity"]
		if updatedReqBody.SentimentSubjectivity.Valid {
			sentimentSubjectivity, isFloatSubjectivity := valueSubjectivity.(float64)
			require.True(t, okSubjectivity, "Expected 'sentiment_subjectivity' field in updated response")
			require.True(t, isFloatSubjectivity, "Expected 'sentiment_subjectivity' value to be a float64 in updated response")
			assert.Equal(t, updatedReqBody.SentimentSubjectivity.Float64, sentimentSubjectivity, "Updated sentiment subjectivity in response does not match")
		} else {
			_, existsSubjectivity := updatedReviewMap["sentiment_subjectivity"]
			assert.False(t, existsSubjectivity, "Did not expect 'sentiment_subjectivity' field in updated response if not valid in request")
		}
	})

	t.Run("update review with non-existent id (expect 404)", func(t *testing.T) {
		updatedReqBody := newTestReview("nonexistent_update")

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(updatedReqBody).
			Put("/api/v1/reviews/99999") // Use non-existent ID in path

		require.Nil(t, err, "Error making PUT request for non-existent id")
		assert.Equal(t, http.StatusNotFound, res.StatusCode(), "Expected 404 Not Found status code")
	})

	t.Run("update review with invalid input (missing required field - expect 400)", func(t *testing.T) {
		invalidReqBody := models.Review{
			// Missing required fields like App, TranslatedReview, Sentiment
			SentimentPolarity: models.NullableFloat64{Float64: 0.1, Valid: true},
		}
		var resBody utils.JSONResponse
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(invalidReqBody).
			SetResult(&resBody).
			Put(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID)) // Use valid ID in path

		require.Nil(t, err, "Error making PUT request with invalid input")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
		assert.Equal(t, "error", resBody.Status, "Expected JSON response status to be 'error'")
	})

	t.Run("update review with invalid id format (expect 400)", func(t *testing.T) {
		updatedReqBody := newTestReview("invalidformat_update")
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(updatedReqBody).
			Put("/api/v1/reviews/abc")

		require.Nil(t, err, "Error making PUT request with invalid id format")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
	})
}

// TestReviewController_DeleteReview tests the DELETE /api/v1/reviews/{id} endpoint
func TestReviewController_DeleteReview(t *testing.T) {
	var createdReviewID int
	reqBody := newTestReview("delete")
	var setupResBody utils.JSONResponse

	// Setup: Create a review to delete
	setupRes, setupErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetResult(&setupResBody).
		Post("/api/v1/reviews")
	require.Nil(t, setupErr)
	require.Equal(t, http.StatusCreated, setupRes.StatusCode, "Setup: Expected 201 Created when creating review for DeleteReview test")

	setupDataMap, ok := setupResBody.Data.(map[string]interface{})
	require.True(t, ok)
	idFloat, ok := setupDataMap["id"].(float64)
	require.True(t, ok)
	createdReviewID = int(idFloat)
	require.NotEqual(t, 0, createdReviewID)

	// No specific cleanup for the review itself needed in this test function's Cleanup,
	// as the deletion is what's being tested.

	t.Run("delete review with valid id", func(t *testing.T) {
		res, err := client.R().Delete(fmt.Sprintf("/api/v1/reviews/%d", createdReviewID))

		require.Nil(t, err, "Error making DELETE request")
		assert.Equal(t, http.StatusOK, res.StatusCode(), "Expected 200 OK status code")
		// Optional: Assert on the success response body if it contains a message
	})

	t.Run("delete review with non-existent id (expect 404)", func(t *testing.T) {
		res, err := client.R().Delete("/api/v1/reviews/99999")

		require.Nil(t, err, "Error making DELETE request for non-existent id")
		assert.Equal(t, http.StatusNotFound, res.StatusCode(), "Expected 404 Not Found status code")
	})

	t.Run("delete review with invalid id format (expect 400)", func(t *testing.T) {
		res, err := client.R().Delete("/api/v1/reviews/abc")

		require.Nil(t, err, "Error making DELETE request with invalid id format")
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), "Expected 400 Bad Request status code")
	})
}
