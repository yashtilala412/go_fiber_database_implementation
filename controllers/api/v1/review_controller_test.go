package v1_test

import (
	"net/http"
	"testing"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/pkg/structs"
	"github.com/stretchr/testify/assert"
)

// TestCreateReview tests the POST /api/v1/reviews endpoint
func TestCreateReview(t *testing.T) {
	t.Run("create review with invalid input", func(t *testing.T) {
		req := structs.Review{
			App: "", // missing required fields
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	t.Run("create review with valid input", func(t *testing.T) {
		req := structs.Review{
			App:                   "MyTestApp",
			TranslatedReview:      "Great app!",
			Sentiment:             "positive",
			SentimentPolarity:     structs.NullableFloat64{Float64: 0.9, Valid: true},
			SentimentSubjectivity: structs.NullableFloat64{Float64: 0.1, Valid: true},
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Post("/api/v1/reviews")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode())
	})

	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM reviews WHERE app = 'MyTestApp'")
		assert.Nil(t, err)
	})
}

// TestGetReviews tests GET /api/v1/reviews
func TestGetReviews(t *testing.T) {
	t.Run("get reviews with valid limit and offset", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/reviews?limit=10&offset=0")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	t.Run("get reviews with invalid limit", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/reviews?limit=-5&offset=0")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	t.Run("get reviews with invalid offset", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/reviews?limit=10&offset=-1")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
}

// TestGetReviewByID tests GET /api/v1/reviews/{id}
func TestGetReviewByID(t *testing.T) {
	t.Run("get review by valid ID", func(t *testing.T) {
		// Assuming ID 1 exists for testing
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/reviews/1")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	t.Run("get review by non-existing ID", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/reviews/99999")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}

// TestUpdateReview tests PUT /api/v1/reviews/{id}
func TestUpdateReview(t *testing.T) {
	t.Run("update review with valid data", func(t *testing.T) {
		req := structs.Review{
			App:                   "UpdatedApp",
			TranslatedReview:      "Updated review text",
			Sentiment:             "neutral",
			SentimentPolarity:     structs.NullableFloat64{Float64: 0.0, Valid: true},
			SentimentSubjectivity: structs.NullableFloat64{Float64: 0.5, Valid: true},
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/reviews/2") // Assume review ID 1 exists

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	t.Run("update review missing required field", func(t *testing.T) {
		req := structs.Review{
			App:              "", // missing required app name
			TranslatedReview: "Review text",
			Sentiment:        "neutral",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/reviews/1")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	t.Run("update review with non-existing ID", func(t *testing.T) {
		req := structs.Review{
			App:                   "UpdatedApp",
			TranslatedReview:      "Updated review text",
			Sentiment:             "neutral",
			SentimentPolarity:     structs.NullableFloat64{Float64: 0.0, Valid: true},
			SentimentSubjectivity: structs.NullableFloat64{Float64: 0.5, Valid: true},
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/reviews/99999")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}

// TestDeleteReview tests DELETE /api/v1/reviews/{id}
func TestDeleteReview(t *testing.T) {
	t.Run("delete review by valid ID", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Delete("/api/v1/reviews/2")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})

	t.Run("delete review by non-existing ID", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Delete("/api/v1/reviews/99999")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode())
	})
}
