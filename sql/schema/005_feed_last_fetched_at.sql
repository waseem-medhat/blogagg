-- +goose Up
ALTER TABLE FEEDS
    ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE FEEDS
    DROP COLUMN last_fetched_at;
