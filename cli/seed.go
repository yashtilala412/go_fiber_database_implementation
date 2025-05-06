package cli

import (
	"database/sql"
	"fmt"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/database"
	_ "github.com/lib/pq" // for postgres dialect
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetSeedCommandDef initializes the seed command
func GetSeedCommandDef(cfg config.AppConfig) cobra.Command {
	seedCmd := cobra.Command{
		Use:   "seed",
		Short: "Seed database with initial data",
		Long:  `This command reads data from CSV files and populates the database tables.`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbConnGoqu, err := database.Connect(cfg.DB)
			if err != nil {
				return fmt.Errorf("failed to connect to database for seeding: %w", err)
			}

			// Get the underlying *sql.DB for closing
			if sqlDB, ok := dbConnGoqu.Db.(*sql.DB); ok {
				defer func() {
					if err := sqlDB.Close(); err != nil {
						fmt.Println("Error closing database connection:", err)
					}
				}()
			} else {
				fmt.Println("Warning: Could not access the underlying *sql.DB to close the connection.")
			}

			err = database.SeedData(cfg, dbConnGoqu, &zap.Logger{})
			if err != nil {
				return fmt.Errorf("failed to seed data: %w", err)
			}
			fmt.Println("Database seeding completed successfully!")
			return nil
		},
	}
	return seedCmd
}
