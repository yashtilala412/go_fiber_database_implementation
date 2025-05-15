package v1_test

import (
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
			Get(fmt.Sprintf("/api/v1/reviews/1", insertedID))

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
