package database

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

// SeedData reads data from CSV files and inserts it into the database.
func SeedData(cfg config.AppConfig, db *goqu.Database, logger *zap.Logger) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer rollback, this will happen if any error occurs.
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic occurred during seeding, rolling back transaction", zap.Any("error", r))
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("Failed to rollback transaction after panic", zap.Error(rbErr))
			}
			panic(r) // Re-panic to propagate.
		} else if err != nil {
			// Rollback only if err is not nil
			logger.Error("Error occurred during seeding, rolling back transaction", zap.Error(err))
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("Failed to rollback transaction after error", zap.Error(rbErr))
			}
		} else {
			// Commit only if there was no error
			if cErr := tx.Commit(); cErr != nil {
				logger.Error("Failed to commit transaction", zap.Error(cErr))
				err = cErr //important: set the error
			}
		}
	}()

	if err := seedAppData(cfg.AppDataCSVPath, tx, logger); err != nil {
		return err // Return the error from seedAppData
	}

	if err := seedReviewData(cfg.ReviewDataCSVPath, tx, logger); err != nil {
		return err // Return the error from seedReviewData
	}
	return nil
}

func seedAppData(csvPath string, tx *goqu.TxDatabase, logger *zap.Logger) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open app_data CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Read the header row and discard it
	if err != nil {
		return fmt.Errorf("failed to read app_data CSV header: %w", err)
	}

	expectedFields := 13

	var appData []map[string]interface{}
	for lineNum := 2; ; lineNum++ {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Warning: Error reading app_data CSV row at line %d: %v\n", lineNum, err)
			continue // Skip rows with read errors other than EOF
		}

		if len(row) != expectedFields {
			fmt.Printf("Warning: Skipping app_data row at line %d with %d fields (expected %d): %v\n", lineNum, len(row), expectedFields, row)
			continue // Skip the current row if the number of fields is wrong
		}

		var rating float64 // Change to float64, not pointer
		if row[2] != "NaN" {
			rating, err = strconv.ParseFloat(row[2], 64)
			if err != nil {
				fmt.Printf("Warning: Skipping row at line %d due to error parsing Rating: %v - Row: %v\n", lineNum, err, row)
				continue
			}
		} else {
			rating = 0.0 // Default value for NaN
		}

		reviews, err := strconv.Atoi(strings.ReplaceAll(row[3], ",", ""))
		if err != nil {
			fmt.Printf("Warning: Skipping row at line %d due to error parsing Reviews: %v - Row: %v\n", lineNum, err, row)
			continue
		}

		data := map[string]interface{}{
			"app":            row[0],
			"category":       row[1],
			"rating":         rating,
			"reviews":        reviews,
			"size":           row[4],
			"installs":       row[5],
			"type":           row[6],
			"price":          row[7],
			"content_rating": row[8],
			"genres":         row[9],
			"last_updated":   row[10],
			"current_ver":    row[11],
			"android_ver":    row[12],
		}
		appData = append(appData, data)
		fmt.Printf("Debug: Inserting row at line %d: %+v\n", lineNum, data) // Debug log
	}

	_, err = tx.Insert("apps").Rows(appData).Executor().Exec()
	if err != nil {
		return fmt.Errorf("failed to insert app_data: %w", err)
	}
	return nil
}

func seedReviewData(csvPath string, tx *goqu.TxDatabase, logger *zap.Logger) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open review_data CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Read the header row
	if err != nil {
		return fmt.Errorf("failed to read review_data CSV header: %w", err)
	}

	expectedFields := 5

	var reviewData []map[string]interface{}
	for lineNum := 2; ; lineNum++ {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Warning: Error reading review_data CSV row at line %d: %v\n", lineNum, err)
			continue
		}

		if len(row) != expectedFields {
			fmt.Printf("Warning: Skipping review_data row at line %d with %d fields (expected %d): %v\n", lineNum, err, row)
			continue
		}

		// Trim spaces from all row values
		for i := range row {
			row[i] = strings.TrimSpace(row[i])
		}

		var sentimentPolarity interface{}
		if row[3] != "nan" && row[3] != "NaN" { // Handle both "nan" and "NaN"
			sentimentPolarity, err = strconv.ParseFloat(row[3], 64)
			if err != nil {
				fmt.Printf("Warning: Skipping row at line %d due to error parsing Sentiment Polarity: %v, Row: %v\n", lineNum, err, row)
				continue
			}
		} else {
			sentimentPolarity = nil // Use nil for NULL
		}

		var sentimentSubjectivity interface{}
		if row[4] != "nan" && row[4] != "NaN" { // Handle both "nan" and "NaN"
			sentimentSubjectivity, err = strconv.ParseFloat(row[4], 64)
			if err != nil {
				fmt.Printf("Warning: Skipping row at line %d due to error parsing Sentiment Subjectivity: %v, Row: %v\n", lineNum, err, row)
				continue
			}
		} else {
			sentimentSubjectivity = nil // Use nil for NULL
		}

		data := map[string]interface{}{
			"app":                    handleString(row[0]),
			"translated_review":      handleString(row[1]), // Changed to snake case
			"sentiment":              handleString(row[2]),
			"sentiment_polarity":     sentimentPolarity,     // Changed to snake case
			"sentiment_subjectivity": sentimentSubjectivity, // Changed to snake case
		}
		reviewData = append(reviewData, data)
		fmt.Printf("Debug: Inserting review row at line %d: %v\n", lineNum, data) // Log the data
	}

	_, err = tx.Insert("reviews").Rows(reviewData).Executor().Exec()
	if err != nil {
		return fmt.Errorf("failed to insert review_data: %w", err)
	}

	return nil
}

// handleString replaces "NaN" or empty strings with a default string value.
func handleString(s string) interface{} { // Return interface{}
	if s == "NaN" || s == "nan" || strings.TrimSpace(s) == "" { // Also check "nan"
		return ""
	}
	return s
}
