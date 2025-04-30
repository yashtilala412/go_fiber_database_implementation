-- +migrate Up

CREATE TABLE review_data (
    App TEXT NOT NULL,
    Translated_Review TEXT NOT NULL,
    Sentiment TEXT NOT NULL,
    Sentiment_Polarity REAL NOT NULL,
    Sentiment_Subjectivity REAL NOT NULL
);
