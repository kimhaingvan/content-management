-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "address" (
    "id" int8 NOT NULL PRIMARY KEY UNIQUE,
    "country" text  DEFAULT 'VN'::text,
    "province_code" text ,
    "province" text ,
    "district_code" text ,
    "district" text,
    "ward" text ,
    "ward_code" text ,
    "address1" text ,
    "created_at" timestamptz(6) NOT NULL,
    "updated_at" timestamptz(6) NOT NULL,
    "deleted_at" timestamptz(6),
    "city" text ,
    "zip" text,
    "address2" text ,
    "note" TEXT
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "address";
-- +goose StatementEnd
