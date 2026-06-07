-- +goose Up
-- above comment required for goose to work
CREATE TABLE users (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email text unique not null
);

-- +goose Down
-- above comment required for goose to work
DROP TABLE users;