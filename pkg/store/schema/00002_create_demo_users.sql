-- +goose Up
-- +goose StatementBegin
INSERT INTO Users (name) VALUES ('Alice');
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO Users (name) VALUES ('Bob');
-- +goose StatementEnd
-- +goose StatementBegin
INSERT INTO Users (name) VALUES ('Carol');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM Users WHERE id = 1;
-- +goose StatementEnd
-- +goose StatementBegin
DELETE FROM Users WHERE id = 2;
-- +goose StatementEnd
-- +goose StatementBegin
DELETE FROM Users WHERE id = 3;
-- +goose StatementEnd
