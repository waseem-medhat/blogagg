-- +goose Up
CREATE TABLE follows (
    id UUID PRIMARY KEY,
    feed_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT user_id_fk FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT feed_id_fk FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT constr_unique_user_feed UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE follows;
