package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

// AppTable represent table name
const AppTable = "apps"

// App model
type App struct {
	AppId         int     `json:"id" db:"id"`
	App           string  `json:"app" db:"app" validate:"required"`
	Category      string  `json:"category" db:"category" validate:"required"`
	Rating        float64 `json:"rating" db:"rating" validate:"required"`
	Reviews       int     `json:"reviews" db:"reviews" validate:"required"`
	Size          string  `json:"size" db:"size" validate:"required"`
	Installs      string  `json:"installs" db:"installs" validate:"required"`
	Type          string  `json:"type" db:"type" validate:"required"`
	Price         string  `json:"price" db:"price" validate:"required"`
	ContentRating string  `json:"content_rating" db:"content_rating" validate:"required"`
	Genres        string  `json:"genres" db:"genres" validate:"required"`
	LastUpdated   string  `json:"last_updated" db:"last_updated" validate:"required"`
	CurrentVer    string  `json:"current_ver" db:"current_ver" validate:"required"`
	AndroidVer    string  `json:"android_ver" db:"android_ver" validate:"required"`
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
func (model *AppModel) GetAppById(id int) (App, error) {
	app := App{}
	found, err := model.db.From(AppTable).Where(goqu.Ex{
		"id": id,
	}).ScanStruct(&app)

	if err != nil {
		return app, err
	}

	if !found {
		return app, sql.ErrNoRows
	}

	return app, nil
}

// InsertApps inserts a new app into the database.
// For AppModel with database-generated ID (SERIAL)
// InsertApps inserts a new app into the database.
func (model *AppModel) InsertApps(app App) (App, error) {
	var insertedID int64

	_, err := model.db.Insert(AppTable).
		Rows(goqu.Record{
			"app":            app.App,
			"category":       app.Category,
			"rating":         app.Rating,
			"reviews":        app.Reviews,
			"size":           app.Size,
			"installs":       app.Installs,
			"type":           app.Type,
			"price":          app.Price,
			"content_rating": app.ContentRating,
			"genres":         app.Genres,
			"last_updated":   app.LastUpdated,
			"current_ver":    app.CurrentVer,
			"android_ver":    app.AndroidVer,
		}).
		Returning("id"). // This makes PostgreSQL return the inserted ID
		Executor().
		ScanVal(&insertedID)

	if err != nil {
		return App{}, err
	}

	app.AppId = int(insertedID)
	return app, nil
}

func (model *AppModel) DeleteApp(id int) error {
	result, err := model.db.Delete(AppTable).Where(goqu.Ex{
		"id": id,
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

func (model *AppModel) UpdateApp(id int, app App) (App, error) {
	result, err := model.db.Update(AppTable).Set(goqu.Record{
		"app":            app.App,
		"category":       app.Category,
		"rating":         app.Rating,
		"reviews":        app.Reviews,
		"size":           app.Size,
		"installs":       app.Installs,
		"type":           app.Type,
		"price":          app.Price,
		"content_rating": app.ContentRating,
		"genres":         app.Genres,
		"last_updated":   app.LastUpdated,
		"current_ver":    app.CurrentVer,
		"android_ver":    app.AndroidVer,
	}).Where(goqu.Ex{"id": id}).Executor().Exec()
	if err != nil {
		return App{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return App{}, err
	}
	if rowsAffected == 0 {
		return App{}, sql.ErrNoRows
	}

	app.AppId = id
	return app, nil
}
