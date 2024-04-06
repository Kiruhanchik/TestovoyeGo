-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cars (
 id SERIAL PRIMARY KEY,
 regNum VARCHAR(255) NOT NULL UNIQUE,
 mark VARCHAR(255) NOT NULL,
 model VARCHAR(255) NOT NULL,
 year INT,
owner_name VARCHAR(255) NOT NULL, 
owner_surname VARCHAR(255) NOT NULL, 
owner_patronymic  VARCHAR(255) NULL
);
-- +goose StatementEnd

CREATE INDEX ON cars (regNum);
CREATE INDEX ON cars (owner_surname);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cars;
-- +goose StatementEnd
