package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// ReviewTable represent table name
const ReviewTable = "review_data"

// Review model
type Review struct {
	ReviewID              int             `json:"review_id" db:"review_id"`
	App                   string          `json:"app" db:"app" validate:"required"`
	TranslatedReview      string          `json:"translated_review" db:"translated_review" validate:"required"`
	Sentiment             string          `json:"sentiment" db:"sentiment" validate:"required"`
	SentimentPolarity     sql.NullFloat64 `json:"sentiment_polarity" db:"sentiment_polarity"`
	SentimentSubjectivity sql.NullFloat64 `json:"sentiment_subjectivity" db:"sentiment_subjectivity"`
}

// ReviewModel implements review related database operations
type ReviewModel struct {
	db *goqu.Database
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
func (model *ReviewModel) GetById(id int) (Review, error) {
	review := Review{}
	found, err := model.db.From(ReviewTable).Where(goqu.Ex{
		"review_id": id,
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
			// Add other unique fields if necessary to ensure correct retrieval
		}).
		ScanStruct(&insertedReview)

	if err != nil {
		return review, err
	}

	if !found {
		return review, sql.ErrNoRows // Or a custom error
	}
	return insertedReview, nil
}

// DeleteApp deletes a review by its ID.
func (model *ReviewModel) DeleteApp(id int) error {
	result, err := model.db.Delete(ReviewTable).Where(goqu.Ex{
		"review_id": id,
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

// UpdateApp updates an existing review by its ID.
func (model *ReviewModel) UpdateApp(id int, review Review) (Review, error) {
	//  Use a transaction to ensure data consistency.
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
		"review_id": id,
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

	// Retrieve the updated record to return it.
	updatedReview := Review{}
	found, err := tx.From(ReviewTable).Where(goqu.Ex{
		"review_id": id,
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
