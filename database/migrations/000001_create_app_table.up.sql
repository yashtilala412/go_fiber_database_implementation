-- +migrate Up

CREATE TABLE apps (
    id SERIAL PRIMARY KEY,
    app TEXT NOT NULL,
    category TEXT NOT NULL,
    rating REAL,
    reviews INTEGER NOT NULL,
    size TEXT NOT NULL,
    installs TEXT NOT NULL,
    type TEXT NOT NULL,
    price TEXT NOT NULL,
    content_rating TEXT NOT NULL,
    genres TEXT NOT NULL,
    last_updated TEXT NOT NULL,
    current_ver TEXT NOT NULL,
    android_ver TEXT NOT NULL
);