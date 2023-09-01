-- +goose Up
-- +goose StatementBegin
CREATE TABLE "authors" (
    "id" integer PRIMARY KEY,
    "name" varchar(255) UNIQUE NOT NULL,
    "email" varchar(255) UNIQUE NOT NULL,
    "password" char(60) NOT NULL
);

ALTER TABLE "todos"
ADD COLUMN "author_id" integer
REFERENCES "authors"("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "authors";

ALTER TABLE "todos"
DROP COLUMN "author_id";
-- +goose StatementEnd
