package v1

import (
	"database/sql"
	"net/http"
	"strconv"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/doug-martin/goqu/v9"
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

func (ac *AppController) GetApp(c *fiber.Ctx) error {
	appID, err := c.ParamsInt(constants.ParamAppID) // Use c.ParamsInt
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidAppID)
	}

	app, err := ac.appService.GetById(appID)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrorAppNotFound)
		}
		ac.logger.Error("error while get app by id", zap.Int("id", appID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.FailedToGetApp)
	}
	return utils.JSONSuccess(c, http.StatusOK, app)
}

func (ac *AppController) GetApps(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(constants.DefaultLimit))) // Use constants
	if err != nil {
		ac.logger.Error("Invalid limit parameter", zap.String("limit", c.Query("limit")), zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrorInvalidLimit)
	}

	offset, err := strconv.Atoi(c.Query("offset", strconv.Itoa(constants.DefaultOffset))) // Use constants
	if err != nil {
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
