package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants" // Import your constants
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ReviewController handles API requests related to review data.
type ReviewController struct {
	reviewService *models.ReviewModel
	logger        *zap.Logger
}

// NewReviewController returns a new ReviewController
func NewReviewController(goqu *goqu.Database, logger *zap.Logger) (*ReviewController, error) {
	reviewModel, err := models.InitReviewModel(goqu)
	if err != nil {
		return nil, err
	}

	return &ReviewController{
		reviewService: &reviewModel,
		logger:        logger,
	}, nil
}

// GetReviews retrieves a paginated list of reviews.
//
//	@Summary		Get Reviews
//	@Description	Fetches reviews from the database with pagination.
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			limit	query	int	false	"Number of reviews to return"
//	@Param			offset	query	int	false	"Offset for pagination"
//	@Success		200	{array}	models.Review
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/reviews [get]

func (rc *ReviewController) GetReviews(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(constants.DefaultLimit)))
	if err != nil {
		rc.logger.Error("Invalid limit parameter", zap.String("limit", c.Query("limit")), zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidLimit)
	}

	offset, err := strconv.Atoi(c.Query("offset", strconv.Itoa(constants.DefaultOffset)))
	if err != nil {
		rc.logger.Error("Invalid offset parameter", zap.String("offset", c.Query("offset")), zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidOffset)
	}

	reviews, err := rc.reviewService.GetReviews(limit, offset)
	if err != nil {
		rc.logger.Error("Failed to get reviews", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.FailedToGetReviews)
	}

	return utils.JSONSuccess(c, http.StatusOK, reviews)
}

// GetReview retrieves a review by ID.
//
//	@Summary		Get Review
//	@Description	Fetches a review using its ID.
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Review ID"
//	@Success		200	{object}	models.Review
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		404	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/reviews/{id} [get]

func (rc *ReviewController) GetReview(c *fiber.Ctx) error {
	reviewID, err := c.ParamsInt(constants.ParamReviewID) //  c.ParamsInt
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidReviewID)
	}

	review, err := rc.reviewService.GetById(reviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrorReviewNotFound)
		}
		rc.logger.Error("error while get review by id", zap.Int("id", reviewID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.FailedToGetReview)
	}
	return utils.JSONSuccess(c, http.StatusOK, review)
}

// CreateReviewData adds a new review.
//
//	@Summary		Create Review
//	@Description	Creates a new review in the database.
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			review	body	models.Review	true	"Review data to create"
//	@Success		201	{object}	models.Review
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/reviews [post]

func (rc *ReviewController) CreateReviewData(c *fiber.Ctx) error {
	var reviewReq models.Review

	err := json.Unmarshal(c.Body(), &reviewReq)
	if err != nil {
		rc.logger.Error("Error unmarshalling request body", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidRequestBody)
	}

	validate := validator.New()
	err = validate.Struct(reviewReq)
	if err != nil {
		rc.logger.Error("Validation error", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	reviewToInsert := models.Review{
		App:                   reviewReq.App,
		TranslatedReview:      reviewReq.TranslatedReview,
		Sentiment:             reviewReq.Sentiment,
		SentimentPolarity:     reviewReq.SentimentPolarity,
		SentimentSubjectivity: reviewReq.SentimentSubjectivity,
	}

	insertedReview, err := rc.reviewService.InsertReviews(reviewToInsert)
	if err != nil {
		rc.logger.Error("Error inserting review data", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrorFiledToCreateReviewApp)
	}

	return utils.JSONSuccess(c, http.StatusCreated, insertedReview)
}

// DeleteReview deletes a review by ID.
//
//	@Summary		Delete Review
//	@Description	Deletes a review by its ID.
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Review ID"
//	@Success		200	{object}	utils.JSONResponse
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		404	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/reviews/{id} [delete]

func (rc *ReviewController) DeleteReview(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params(constants.ParamReviewID))
	if err != nil {
		rc.logger.Error("Error parsing review ID", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidReviewID)
	}

	err = rc.reviewService.DeleteApp(id)
	if err != nil {
		if err == sql.ErrNoRows {
			rc.logger.Warn("Review not found", zap.Int("id", id))
			return utils.JSONError(c, http.StatusNotFound, constants.ErrorReviewNotFound)
		}
		rc.logger.Error("Error deleting review", zap.Error(err), zap.Int("id", id))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrorFaiedToDeleteReview)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.ReviewsDeletedSuccessfully)
}

// UpdateReview updates a review.
//
//	@Summary		Update Review
//	@Description	Updates an existing review in the database.
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param		id		path		int			true	"Review ID"
//	@Param		review	body		models.Review	true	"Updated review data"
//	@Success		200	{object}	models.Review
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		404	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/reviews/{id} [put]
func (rc *ReviewController) UpdateReview(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params(constants.ParamReviewID))
	if err != nil {
		rc.logger.Error("Error parsing review ID", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidReviewID)
	}

	var updatedReview models.Review
	if err := json.Unmarshal(c.Body(), &updatedReview); err != nil {
		rc.logger.Error("Error parsing request body", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidRequestBody)
	}

	// Validate the request body.
	validate := validator.New()
	err = validate.Struct(updatedReview)
	if err != nil {
		rc.logger.Error("Validation error", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, utils.ValidatorErrorString(err)) //  Adapt this as needed
	}

	updatedReview, err = rc.reviewService.UpdateApp(id, updatedReview)
	if err != nil {
		if err == sql.ErrNoRows {
			rc.logger.Warn("Review not found", zap.Int("id", id))
			return utils.JSONError(c, http.StatusNotFound, constants.ErrorReviewNotFound)
		}
		rc.logger.Error("Error updating review", zap.Error(err), zap.Int("id", id))
		return utils.JSONError(c, http.StatusInternalServerError, constants.FailedToUpdateReviews)
	}

	return utils.JSONSuccess(c, http.StatusOK, updatedReview)
}
