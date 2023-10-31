-- +goose Up
-- +goose StatementBegin
CREATE DATABASE spiredemo;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE spiredemo;
-- +goose StatementEnd
