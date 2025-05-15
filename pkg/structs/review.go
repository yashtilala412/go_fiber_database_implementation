package structs

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Review defines the review structure for API responses
type Review struct {
	ID                    int             `json:"id" db:"id"`
	App                   string          `json:"app" db:"app" validate:"required"`
	TranslatedReview      string          `json:"translated_review" db:"translated_review" validate:"required"`
	Sentiment             string          `json:"sentiment" db:"sentiment" validate:"required"`
	SentimentPolarity     NullableFloat64 `json:"sentiment_polarity" db:"sentiment_polarity"`
	SentimentSubjectivity NullableFloat64 `json:"sentiment_subjectivity" db:"sentiment_subjectivity"`
}

// ReqCreateReview defines the request body for creating a review
type ReqCreateReview struct {
	App                   string          `json:"app" db:"app" validate:"required"`
	TranslatedReview      string          `json:"translated_review" db:"translated_review" validate:"required"`
	Sentiment             string          `json:"sentiment" db:"sentiment" validate:"required"`
	SentimentPolarity     NullableFloat64 `json:"sentiment_polarity" db:"sentiment_polarity"`
	SentimentSubjectivity NullableFloat64 `json:"sentiment_subjectivity" db:"sentiment_subjectivity"`
}

// NullableFloat64 is a custom type that handles nullable float64 values
type NullableFloat64 struct {
	Float64 float64 `json:"value" swaggertype:"primitive,number"`
	Valid   bool    `json:"valid" swaggertype:"primitive,boolean"`
}

// MarshalJSON implements json.Marshaler interface
func (nf NullableFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (nf *NullableFloat64) UnmarshalJSON(data []byte) error {
	// If the input is null, set Valid to false
	if string(data) == "null" {
		nf.Valid = false
		return nil
	}

	// Otherwise try to unmarshal into float64
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	nf.Float64 = f
	nf.Valid = true
	return nil
}

// Scan implements the sql.Scanner interface
func (nf *NullableFloat64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil {
		return err
	}

	nf.Float64 = f.Float64
	nf.Valid = f.Valid
	return nil
}

// Value implements the driver.Valuer interface
func (nf NullableFloat64) Value() (driver.Value, error) {
	if !nf.Valid {
		return nil, nil
	}
	return nf.Float64, nil
}

// String returns a string representation of the value
func (nf NullableFloat64) String() string {
	if !nf.Valid {
		return "NULL"
	}
	return fmt.Sprintf("%f", nf.Float64)
}
