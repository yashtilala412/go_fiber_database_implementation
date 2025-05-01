package cli

import (
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	migrationCmd := GetMigrationCommandDef(cfg)
	apiCmd := GetAPICommandDef(cfg, logger)
	seedCmd := GetSeedCommandDef(cfg) // Add the seed command

	rootCmd := &cobra.Command{Use: "golang-api"}
	rootCmd.AddCommand(&migrationCmd, &apiCmd, &seedCmd)
	return rootCmd.Execute()
}
