package v1_test

import (
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
