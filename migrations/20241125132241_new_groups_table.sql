-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS groups (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS groups;

-- +goose StatementEnd
