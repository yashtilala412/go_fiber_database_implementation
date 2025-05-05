package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// AppTable represent table name
const AppTable = "app_data"

// App model
type App struct {
	AppID         int     `json:"app_id" db:"app_id"`
	App           string  `json:"app" db:"app" validate:"required"`
	Category      string  `json:"category" db:"category" validate:"required"`
	Rating        float64 `json:"rating" db:"rating"`
	Reviews       int     `json:"reviews" db:"reviews"`
	Size          string  `json:"size" db:"size"`
	Installs      string  `json:"installs" db:"installs"`
	Type          string  `json:"type" db:"type"`
	Price         string  `json:"price" db:"price"`
	ContentRating string  `json:"content_rating" db:"content_rating"` // Changed to snake case
	Genres        string  `json:"genres" db:"genres"`
	LastUpdated   string  `json:"last_updated" db:"last_updated"` // Changed to snake case
	CurrentVer    string  `json:"current_ver" db:"current_ver"`   // Changed to snake case
	AndroidVer    string  `json:"android_ver" db:"android_ver"`   // Changed to snake case
}

// AppModel implements app related database operations
type AppModel struct {
	db *goqu.Database
}

// InitAppModel Init model
func InitAppModel(goqu *goqu.Database) (AppModel, error) {
	return AppModel{
		db: goqu,
	}, nil
}

// GetApps lists all apps.
func (model *AppModel) GetApps(limit, offset int) ([]App, error) {
	var apps []App
	query := model.db.From(AppTable)

	if limit > 0 {
		query = query.Limit(uint(limit))
	}
	if offset >= 0 {
		query = query.Offset(uint(offset))
	}

	if err := query.ScanStructs(&apps); err != nil {
		return nil, err
	}
	return apps, nil
}

// GetById gets an app by its ID.  It retrieves all fields from the database.
func (model *AppModel) GetById(id int) (App, error) {
	app := App{}
	found, err := model.db.From(AppTable).Where(goqu.Ex{
		"app_id": id,
	}).ScanStruct(&app)

	if err != nil {
		return app, err
	}

	if !found {
		return app, sql.ErrNoRows
	}

	return app, nil
}

// InsertAppData inserts a new app into the database.
func (model *AppModel) InsertAppData(app App) (App, error) {
	_, err := model.db.Insert(AppTable).Rows(goqu.Record{
		"app":            app.App,
		"category":       app.Category,
		"rating":         app.Rating,
		"reviews":        app.Reviews,
		"size":           app.Size,
		"installs":       app.Installs,
		"type":           app.Type,
		"price":          app.Price,
		"content_rating": app.ContentRating, // Changed to snake case
		"genres":         app.Genres,
		"last_updated":   app.LastUpdated, // Changed to snake case
		"current_ver":    app.CurrentVer,  // Changed to snake case
		"android_ver":    app.AndroidVer,  // Changed to snake case
	}).Executor().Exec()
	if err != nil {
		return app, err
	}

	//  we should query the database to get the complete record, including the generated app_id.
	var insertedApp App
	found, err := model.db.From(AppTable).
		Where(goqu.Ex{ // Assuming other fields are unique enough to identify the inserted row.
			"app":      app.App,
			"category": app.Category,
			// Add other fields to uniquely identify the record.
		}).
		ScanStruct(&insertedApp)

	if err != nil {
		return app, err
	}
	if !found {
		return app, sql.ErrNoRows // Or some other error to indicate that the record was not found.
	}
	return insertedApp, nil // Return the full record.
}
func (model *AppModel) DeleteByID(id int) error {
	result, err := model.db.Delete(AppTable).Where(goqu.Ex{
		"app_id": id,
	}).Executor().Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Return a specific error for "not found"
	}
	return nil
}

// UpdateByID updates an existing app by its ID.
func (model *AppModel) UpdateByID(id int, app App) (App, error) {
	//  Use a transaction to ensure data consistency.
	tx, err := model.db.Begin()
	if err != nil {
		return App{}, err
	}
	defer tx.Rollback() // Rollback if any error occurs

	// Update the record.
	result, err := tx.Update(AppTable).Set(goqu.Record{
		"app":            app.App,
		"category":       app.Category,
		"rating":         app.Rating,
		"reviews":        app.Reviews,
		"size":           app.Size,
		"installs":       app.Installs,
		"type":           app.Type,
		"price":          app.Price,
		"content_rating": app.ContentRating, // Changed to snake case
		"genres":         app.Genres,
		"last_updated":   app.LastUpdated, // Changed to snake case
		"current_ver":    app.CurrentVer,  // Changed to snake case
		"android_ver":    app.AndroidVer,  // Changed to snake case
	}).Where(goqu.Ex{
		"app_id": id,
	}).Executor().Exec()
	if err != nil {
		return App{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return App{}, err
	}
	if rowsAffected == 0 {
		return App{}, sql.ErrNoRows // Return error if no rows were updated
	}

	// Retrieve the updated record to return it.
	updatedApp := App{}
	found, err := tx.From(AppTable).Where(goqu.Ex{
		"app_id": id,
	}).ScanStruct(&updatedApp)
	if err != nil {
		return App{}, err
	}
	if !found {
		return App{}, sql.ErrNoRows //Should not happen, but handle it.
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return App{}, err
	}
	return updatedApp, nil
}
