-- +migrate Up

CREATE TABLE review_data (
    review_id SERIAL PRIMARY KEY,
    app TEXT NOT NULL,
    translated_review TEXT NOT NULL,
    sentiment TEXT NOT NULL,
    sentiment_polarity REAL,        
    sentiment_subjectivity REAL     
);
