-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shop (
    id int8 NOT NULL UNIQUE PRIMARY KEY ,
    name TEXT NOT NULL,
    image_url TEXT,
    phone TEXT,
    email TEXT,
    status int2,
    code TEXT NOT NULL,
    website_url TEXT,
    address_id int8 NOT NULL REFERENCES address(id),
    owner_id int8 NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop;
-- +goose StatementEnd
