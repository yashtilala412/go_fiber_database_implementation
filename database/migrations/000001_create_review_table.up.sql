-- +migrate Up

CREATE TABLE review_data (
    review_id SERIAL PRIMARY KEY,
    App TEXT NOT NULL,
    Translated_Review TEXT NOT NULL,
    Sentiment TEXT NOT NULL,
    Sentiment_Polarity REAL,        
    Sentiment_Subjectivity REAL     
);
