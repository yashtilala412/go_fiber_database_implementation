-- +migrate Up

CREATE TABLE app_data (
    app_id SERIAL PRIMARY KEY,  -- Added app_id as SERIAL PRIMARY KEY
    App TEXT NOT NULL,
    Category TEXT NOT NULL,
    Rating REAL,
    Reviews INTEGER NOT NULL,
    Size TEXT NOT NULL,
    Installs TEXT NOT NULL,
    Type TEXT NOT NULL,
    Price TEXT NOT NULL,
    "Content Rating" TEXT NOT NULL,
    Genres TEXT NOT NULL,
    "Last Updated" TEXT NOT NULL,
    "Current Ver" TEXT NOT NULL,
    "Android Ver" TEXT NOT NULL
);
