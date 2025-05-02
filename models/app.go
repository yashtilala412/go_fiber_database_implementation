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
	ContentRating string  `json:"content_rating" db:"Content Rating"`
	Genres        string  `json:"genres" db:"genres"`
	LastUpdated   string  `json:"last_updated" db:"Last Updated"`
	CurrentVer    string  `json:"current_ver" db:"Current Ver"`
	AndroidVer    string  `json:"android_ver" db:"Android Ver"`
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
