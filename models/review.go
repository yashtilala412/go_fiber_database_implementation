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
