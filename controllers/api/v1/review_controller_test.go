package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestCreateReview(t *testing.T) {
	t.Run("create review with valid input", func(t *testing.T) {
		req := structs.ReqCreateReview{
			App:              "Test App Name update",
			TranslatedReview: "This is a test review. update",
			Sentiment:        "Positive",
			SentimentPolarity: structs.NullableFloat64{
				Float64: 8.5,
				Valid:   true,
			},
			SentimentSubjectivity: structs.NullableFloat64{
				Float64: 0.5,
				Valid:   true,
			},
		}

		res, err := client.
			R().
			SetBody(req).
			SetHeader("Content-Type", "application/json").
			Post("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())
	})

	t.Run("create review missing required fields", func(t *testing.T) {
		req := map[string]interface{}{
			"app":               "",
			"translated_review": "",
			"sentiment":         "",
		}

		res, err := client.
			R().
			SetBody(req).
			SetHeader("Content-Type", "application/json").
			Post("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
	t.Run("create review with valid input but boundary sentiment values", func(t *testing.T) {
		req := structs.ReqCreateReview{
			App:              "EdgeApp",
			TranslatedReview: "Zero polarity and subjectivity test",
			Sentiment:        "Neutral",
			SentimentPolarity: structs.NullableFloat64{
				Float64: 0.0,
				Valid:   true,
			},
			SentimentSubjectivity: structs.NullableFloat64{
				Float64: 0.0,
				Valid:   true,
			},
		}

		res, err := client.
			R().
			SetBody(req).
			SetHeader("Content-Type", "application/json").
			Post("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())
	})
	t.Cleanup(func() {
		_, err := db.Exec(`DELETE FROM reviews WHERE app IN (?, ?)`, "Test App Name update", "EdgeApp")
		assert.Nil(t, err)
	})
}
func TestGetReviewByID(t *testing.T) {
	var insertedID int

	// First create a review to retrieve
	t.Run("create review for fetching", func(t *testing.T) {
		req := structs.ReqCreateReview{
			App:                   "Test App For GetByID",
			TranslatedReview:      "This is a test review for get by id",
			Sentiment:             "positive",
			SentimentPolarity:     structs.NullableFloat64{Float64: 0.7, Valid: true},
			SentimentSubjectivity: structs.NullableFloat64{Float64: 0.6, Valid: true},
		}

		res, err := client.
			R().
			SetBody(req).
			SetResult(&structs.Review{}).
			Post("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())

		createdReview := res.Result().(*structs.Review)
		insertedID = createdReview.ReviewID
	})

	// Test Case: Get review by valid ID
	t.Run("get review by valid id", func(t *testing.T) {
		res, err := client.
			R().
			SetResult(&structs.Review{}).
			Get(fmt.Sprintf("/api/v1/reviews/%d", insertedID))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	// Test Case: Get review by invalid ID
	t.Run("get review by invalid id", func(t *testing.T) {
		res, err := client.
			R().
			Get("/api/v1/reviews/999999") // Assuming this ID doesn't exist

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})

	// Cleanup
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM reviews WHERE id = $1", insertedID)
		assert.Nil(t, err)
	})
}
func TestGetReviews(t *testing.T) {
	// Setup: Create a review using POST endpoint
	setupReview := map[string]interface{}{
		"app":                    "TestAppGet",
		"translated_review":      "Test fetching review",
		"sentiment":              "Positive",
		"sentiment_polarity":     map[string]interface{}{"value": 0.9, "valid": true},
		"sentiment_subjectivity": map[string]interface{}{"value": 0.7, "valid": true},
	}

	resCreate, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(setupReview).
		Post("/api/v1/reviews")

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resCreate.StatusCode())

	t.Run("get all reviews successfully", func(t *testing.T) {
		res, err := client.
			R().
			SetHeader("Accept", "application/json").
			Get("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())

		var response []map[string]interface{}
		err = json.Unmarshal(res.Body(), &response)
		assert.Nil(t, err)
		assert.True(t, len(response) >= 1)
	})

	t.Run("get reviews with negative limit", func(t *testing.T) {
		res, err := client.
			R().
			SetQueryParam("limit", "-10").
			Get("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	t.Run("get reviews with invalid limit format", func(t *testing.T) {
		res, err := client.
			R().
			SetQueryParam("limit", "abc").
			Get("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	// Cleanup
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM reviews WHERE app = 'TestAppGet'")
		assert.Nil(t, err)
	})
}
func TestUpdateReview(t *testing.T) {
	// Step 1: Create a review to update
	setupReview := map[string]interface{}{
		"app":                    "TestAppUpdate",
		"translated_review":      "Initial review text",
		"sentiment":              "Neutral",
		"sentiment_polarity":     map[string]interface{}{"value": 0.0, "valid": true},
		"sentiment_subjectivity": map[string]interface{}{"value": 0.5, "valid": true},
	}

	resCreate, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(setupReview).
		Post("/api/v1/reviews")

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resCreate.StatusCode())

	var created map[string]interface{}
	_ = json.Unmarshal(resCreate.Body(), &created)
	reviewID := int(created["id"].(float64))

	t.Run("update review with valid data", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"app":                    "TestAppUpdate",
			"translated_review":      "Updated review text",
			"sentiment":              "Positive",
			"sentiment_polarity":     0.7,
			"sentiment_subjectivity": 0.9,
		}

		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(updateBody).
			Put(fmt.Sprintf("/api/v1/reviews/%d", reviewID))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	t.Run("update review with invalid ID", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"app":                    "InvalidUpdate",
			"translated_review":      "Bad update",
			"sentiment":              "Neutral",
			"sentiment_polarity":     0.0,
			"sentiment_subjectivity": 0.5,
		}

		res, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetBody(updateBody).
			Put("/api/v1/reviews/999999")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})

	// Cleanup
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM reviews WHERE app = 'TestAppUpdate'")
		assert.Nil(t, err)
	})
}
func TestDeleteReview(t *testing.T) {
	// Step 1: Create review first using helper
	body := map[string]interface{}{
		"app":                    "TestAppDelete",
		"translated_review":      "This is a test delete review",
		"sentiment":              "neutral",
		"sentiment_polarity":     map[string]interface{}{"value": 0.3, "valid": true},
		"sentiment_subjectivity": map[string]interface{}{"value": 0.6, "valid": true},
	}

	resCreate, err := client.
		R().
		SetBody(body).
		SetHeader("Content-Type", "application/json").
		Post("/api/v1/reviews")

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resCreate.StatusCode())

	// Extract ID
	var created map[string]interface{}
	_ = json.Unmarshal(resCreate.Body(), &created)
	reviewID := int(created["id"].(float64))

	// Step 2: Delete the review
	t.Run("delete review with valid ID", func(t *testing.T) {
		res, err := client.
			R().
			Delete(fmt.Sprintf("/api/v1/reviews/%d", reviewID))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	// Step 3: Try deleting again (should fail)
	t.Run("delete review with non-existent ID", func(t *testing.T) {
		res, err := client.
			R().
			Delete(fmt.Sprintf("/api/v1/reviews/%d", reviewID))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}
