package v1

import (
	"context"
	"net/http"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthController struct {
	db     *goqu.Database
	logger *zap.Logger
}

func NewHealthController(db *goqu.Database, logger *zap.Logger) (*HealthController, error) {
	return &HealthController{
		db:     db,
		logger: logger,
	}, nil
}

// Overall check overall health of application as well as dependencies health check
// @Summary Overall health check
// @Description Overall health check of application as well as dependencies health check
// @Tags Healthcheck
// @ID overallHealthCheck
// @Produce json
// @Success 200 {object} utils.JSONResponse "Health check successful"
// @Failure 500 {object} utils.JSONResponse "Internal server error during health check"
// @Router /healthz [get]
func (hc *HealthController) Overall(ctx *fiber.Ctx) error {
	// --- Add this log to see if the handler is being reached ---
	hc.logger.Info("Received request for /healthz")
	// --- End Add ---

	err := healthDb(hc.db) // This checks the database
	if err != nil {
		// --- Add this log to see if the DB check failed ---
		hc.logger.Error("Database health check failed", zap.Error(err))
		// --- End Add ---
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	// --- Add this log to confirm success ---
	hc.logger.Info("Database health check successful")
	// --- End Add ---

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Self health check
// @Summary Self health check
// @Description Basic self health check without dependency checks.
// @Tags Healthcheck
// @ID selfHealthCheck
// @Produce json
// @Success 200 {object} utils.JSONResponse "Health check successful"
// @Router /healthz/self [get]
func (hc *HealthController) Self(ctx *fiber.Ctx) error {
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Database health check
// @Summary Database health check
// @Description Database health check
// @Tags Healthcheck
// @ID dbHealthCheck
// @Produce json
// @Success 200 {object} utils.JSONResponse "Database health check successful"
// @Failure 500 {object} utils.JSONResponse "Internal server error during database health check"
// @Router /healthz/db [get]
func (hc *HealthController) Db(ctx *fiber.Ctx) error {
	err := healthDb(hc.db)
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

///////////////////////
// HealthCheck CORE
//////////////////////

func healthDb(db *goqu.Database) error {
	// Reverted to original implementation
	_, err := db.ExecContext(context.TODO(), "SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
