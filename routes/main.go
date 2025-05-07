package routes

import (
	"fmt"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	controllers "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/controllers/api/v1"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	// Adjust the import path if necessary
	pMetrics "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/pkg/prometheus"
)

// Setup function to include App routes
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error { // Added pMetrics
	router := app.Group("/api")
	v1 := router.Group("/v1")

	// Setup other routes...
	err := setupAppController(v1, goqu, logger, pMetrics) // Pass pMetrics
	if err != nil {
		return err
	}
	// Setup Review routes
	err = setupReviewController(v1, goqu, logger, pMetrics)
	if err != nil {
		return err
	}
	return nil
}

func setupAppController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error { // Added pMetrics
	appController, err := controllers.NewAppController(goqu, logger)
	if err != nil {
		return err
	}

	appRouter := v1.Group("/apps") // Define the /apps route group

	// Define the specific routes within the /apps group
	appRouter.Get(fmt.Sprintf("/:%s", constants.ParamAppID), appController.GetApp) // GET /api/v1/apps/:appId
	appRouter.Get("/", appController.GetApps)
	appRouter.Post("/", appController.CreateApp) // GET /api/v1/apps/
	appRouter.Delete(fmt.Sprintf("/:%s", constants.ParamAppID), appController.DeleteApp)
	appRouter.Put(fmt.Sprintf("/:%s", constants.ParamAppID), appController.UpdateApp)
	return nil
}
func setupReviewController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error {
	reviewController, err := controllers.NewReviewController(goqu, logger)
	if err != nil {
		return err
	}

	reviewRouter := v1.Group("/reviews")

	reviewRouter.Get(fmt.Sprintf("/:%s", constants.ParamReviewID), reviewController.GetReview) // GET /api/v1/reviews/:id
	reviewRouter.Get("/", reviewController.GetReviews)
	reviewRouter.Post("/", reviewController.CreateReviewData)
	reviewRouter.Delete(fmt.Sprintf("/:%s", constants.ParamReviewID), reviewController.DeleteReview)
	reviewRouter.Put(fmt.Sprintf("/:%s", constants.ParamReviewID), reviewController.UpdateReview)

	return nil
}
