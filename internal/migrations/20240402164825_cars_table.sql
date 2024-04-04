-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars (
 id SERIAL PRIMARY KEY,
 regNum VARCHAR(255) NOT NULL,
 mark VARCHAR(255) NOT NULL,
 model VARCHAR(255) NOT NULL,
 owner integer not NULL,
 year integer,
 created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cars;
-- +goose StatementEnd
