package cli

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/database"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	pMetrics "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/pkg/prometheus"
)

// GetAPICommandDef runs app
func GetAPICommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Create fiber app
			app := fiber.New(fiber.Config{})

			promMetrics := pMetrics.InitPrometheusMetrics()

			// Database connection
			db, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}

			// Setup routes
			err = routes.Setup(app, db, logger, promMetrics)
			if err != nil {
				return err
			}

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

			// Start server in a goroutine
			go func() {
				if err := app.Listen(cfg.Host + ":" + cfg.Port); err != nil {
					logger.Panic(err.Error())
				}
			}()

			<-interrupt
			logger.Info("gracefully shutting down...")
			if err := app.Shutdown(); err != nil {
				logger.Panic("error while shutting down server", zap.Error(err))
			}

			logger.Info("server stopped receiving new requests.")
			return nil
		},
	}

	return apiCommand
}
