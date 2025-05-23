package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

// AppController for app controllers
type AppController struct {
	appService *models.AppModel // Use the AppModel directly
	logger     *zap.Logger
}

// NewAppController returns a new AppController
func NewAppController(goqu *goqu.Database, logger *zap.Logger) (*AppController, error) {
	appModel, err := models.InitAppModel(goqu) // Initialize AppModel
	if err != nil {
		return nil, err
	}

	return &AppController{
		appService: &appModel, // Use the initialized AppModel
		logger:     logger,
	}, nil
}

// GetApp retrieves a single app by ID.
//
//	@Summary		Get App
//	@Description	Fetches an app by its ID.
//	@Tags			Apps
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"App ID"
//	@Success		200	{object}	models.App
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		404	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/apps/{id} [get]

func (ac *AppController) GetApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt(constants.ParamAppID) // Use c.ParamsInt
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidAppID)
	}

	app, err := ac.appService.GetAppById(appID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrorAppNotFound)
		}
		ac.logger.Error("error while get app by id", zap.Int("id", appID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.FailedToGetApp)
	}
	return utils.JSONSuccess(c, http.StatusOK, app)
}

// GetApps fetches a list of apps.
//
//	@Summary		Get Apps
//	@Description	Retrieves a paginated list of apps from the database.
//	@Tags			Apps
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Number of records to return"
//	@Param			page	query		int	false	"Page number"
//	@Success		200	{array}	models.App
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/apps [get]

func (ac *AppController) GetApps(c *fiber.Ctx) error {
	const MaxLimit = 500 // Set maximum allowed limit

	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(constants.DefaultLimit)))
	if err != nil || limit <= 0 {
		ac.logger.Error("Invalid limit parameter", zap.String("limit", c.Query("limit")), zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidLimit)
	}

	// Check if limit exceeds MaxLimit
	if limit > MaxLimit {
		ac.logger.Warn("Requested limit exceeds maximum", zap.Int("limit", limit))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorlimitAccess)
	}

	offset, err := strconv.Atoi(c.Query("offset", strconv.Itoa(constants.DefaultOffset)))
	if err != nil || offset < 0 {
		ac.logger.Error("Invalid offset parameter", zap.String("offset", c.Query("offset")), zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidOffset)
	}

	apps, err := ac.appService.GetApps(limit, offset)
	if err != nil {
		ac.logger.Error("Failed to get apps", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.FailedToGetApp)
	}

	return utils.JSONSuccess(c, http.StatusOK, apps)
}

// CreateApp creates a new app.
//
//	@Summary		Create App
//	@Description	Creates a new application entry in the database.
//	@Tags			Apps
//	@Accept			json
//	@Produce		json
//	@Param			app	body		models.App	true	"App data to create"
//	@Success		201	{object}	models.App
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/apps [post]

func (ac *AppController) CreateApp(c *fiber.Ctx) error {
	var appReq models.App // Use the App struct from your models

	// Parse the request body into the App struct.
	err := json.Unmarshal(c.Body(), &appReq)
	if err != nil {
		ac.logger.Error("Error unmarshalling request body", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidRequestBody+err.Error())
	}

	// Validate the request body.
	validate := validator.New()
	err = validate.Struct(appReq)
	if err != nil {
		ac.logger.Error("Validation error", zap.Error(err))
		validationErrors := utils.ValidatorErrorString(err)
		return utils.JSONError(c, http.StatusBadRequest, validationErrors) //  Adapt this as needed.  Send the validation errors.
	}

	// Insert the app data into the database.
	insertedApp, err := ac.appService.InsertApps(appReq)
	if err != nil {
		ac.logger.Error("Error inserting app data", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrorFiledToCreateApp) //Use a constant
	}

	// Return the newly created app data, including the generated ID.
	return utils.JSONSuccess(c, http.StatusCreated, insertedApp)
}

// DeleteApp deletes an app by ID.
//
//	@Summary		Delete App
//	@Description	Removes an app from the database.
//	@Tags			Apps
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"App ID"
//	@Success		200	{object}	utils.JSONResponse
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		404	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/apps/{id} [delete]

func (ac *AppController) DeleteApp(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params(constants.ParamAppID))
	if err != nil {
		ac.logger.Error("Error parsing app ID", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidAppID)
	}

	err = ac.appService.DeleteApp(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ac.logger.Warn("App not found", zap.Int("id", id))
			return utils.JSONError(c, http.StatusNotFound, constants.ErrorAppNotFound)
		}
		ac.logger.Error("Error deleting app", zap.Error(err), zap.Int("id", id))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrorFaiedToDeleteApp)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AppsDeletedSuccessfully)
}

// UpdateApp updates an existing app.
//
//	@Summary		Update App
//	@Description	Updates app data in the database.
//	@Tags			Apps
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int			true	"App ID"
//	@Param			app		body	models.App	true	"Updated app data"
//	@Success		200	{object}	models.App
//	@Failure		400	{object}	utils.JSONResponse
//	@Failure		404	{object}	utils.JSONResponse
//	@Failure		500	{object}	utils.JSONResponse
//	@Router			/api/v1/apps/{id} [put]

func (ac *AppController) UpdateApp(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params(constants.ParamAppID))
	if err != nil {
		ac.logger.Error("Error parsing app ID", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidAppID)
	}

	var updatedApp models.App
	if err := c.BodyParser(&updatedApp); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidRequestBody)
	}

	// Validate the request body.
	validate := validator.New()
	err = validate.Struct(updatedApp)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, utils.ValidatorErrorString(err)) //  Adapt this as needed
	}

	updatedApp, err = ac.appService.UpdateApp(id, updatedApp)
	if err != nil {
		if err == sql.ErrNoRows {
			ac.logger.Warn("App not found", zap.Int("id", id))
			return utils.JSONError(c, http.StatusNotFound, constants.ErrorAppNotFound)
		}
		ac.logger.Error("Error updating app", zap.Error(err), zap.Int("id", id))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrorFiledToUpdateApp)
	}

	return utils.JSONSuccess(c, http.StatusOK, updatedApp)
}
