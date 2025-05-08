package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

// ReviewTable represent table name
const ReviewTable = "reviews"

// Review model
type Review struct {
	ReviewID              int             `json:"id" db:"id"`
	App                   string          `json:"app" db:"app" validate:"required"`
	TranslatedReview      string          `json:"translated_review" db:"translated_review" validate:"required"`
	Sentiment             string          `json:"sentiment" db:"sentiment" validate:"required"`
	SentimentPolarity     NullableFloat64 `json:"sentiment_polarity" db:"sentiment_polarity"`
	SentimentSubjectivity NullableFloat64 `json:"sentiment_subjectivity" db:"sentiment_subjectivity"`
}

// ReviewModel implements review related database operations
type ReviewModel struct {
	db *goqu.Database
}

// because swagger does not handle sql.Nullfloat64 directly
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

// InitReviewModel Init model
func InitReviewModel(goqu *goqu.Database) (ReviewModel, error) {
	return ReviewModel{
		db: goqu,
	}, nil
}

// GetReviews lists all reviews.
func (model *ReviewModel) GetReviews(limit, offset int) ([]Review, error) {
	var reviews []Review
	query := model.db.From(ReviewTable)

	if limit > 0 {
		query = query.Limit(uint(limit))
	}
	if offset >= 0 {
		query = query.Offset(uint(offset))
	}

	if err := query.ScanStructs(&reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetById gets a review by its ID.
func (model *ReviewModel) GetReviewById(id int) (Review, error) {
	review := Review{}
	found, err := model.db.From(ReviewTable).Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&review)

	if err != nil {
		return review, err
	}

	if !found {
		return review, sql.ErrNoRows
	}

	return review, nil
}

// InsertReviews inserts a new review into the database.
func (model *ReviewModel) InsertReviews(review Review) (Review, error) {
	_, err := model.db.Insert(ReviewTable).Rows(goqu.Record{
		"app":                    review.App,
		"translated_review":      review.TranslatedReview,
		"sentiment":              review.Sentiment,
		"sentiment_polarity":     review.SentimentPolarity,
		"sentiment_subjectivity": review.SentimentSubjectivity,
	}).Executor().Exec()
	if err != nil {
		return review, err
	}

	// Retrieve the inserted review to get the generated ID.
	var insertedReview Review
	found, err := model.db.From(ReviewTable).
		Where(goqu.Ex{
			"app":               review.App,
			"translated_review": review.TranslatedReview,
		}).
		ScanStruct(&insertedReview)

	if err != nil {
		return review, err
	}

	if !found {
		return review, sql.ErrNoRows
	}
	return insertedReview, nil
}

// DeleteApp deletes a review by its ID.
func (model *ReviewModel) DeleteApp(id int) error {
	result, err := model.db.Delete(ReviewTable).Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Return a specific error for "not found"
	}
	return nil
}

func (model *ReviewModel) UpdateApp(id int, review Review) (Review, error) {
	tx, err := model.db.Begin()
	if err != nil {
		return Review{}, err
	}
	defer tx.Rollback() // Rollback if any error occurs

	// Update the record.
	result, err := tx.Update(ReviewTable).Set(goqu.Record{
		"app":                    review.App,
		"translated_review":      review.TranslatedReview,
		"sentiment":              review.Sentiment,
		"sentiment_polarity":     review.SentimentPolarity,
		"sentiment_subjectivity": review.SentimentSubjectivity,
	}).Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err != nil {
		return Review{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Review{}, err
	}
	if rowsAffected == 0 {
		return Review{}, sql.ErrNoRows // Return error if no rows were updated
	}

	updatedReview := Review{}
	found, err := tx.From(ReviewTable).Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&updatedReview)
	if err != nil {
		return Review{}, err
	}
	if !found {
		return Review{}, sql.ErrNoRows //Should not happen, but handle it.
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return Review{}, err
	}
	return updatedReview, nil
}
