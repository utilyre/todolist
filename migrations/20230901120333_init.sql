-- +goose Up
-- +goose StatementBegin
CREATE TABLE "todos" (
    "id" integer PRIMARY KEY,
    "title" varchar(16),
    "body" varchar(1024)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "todos";
-- +goose StatementEnd
