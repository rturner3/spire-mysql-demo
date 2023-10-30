-- +goose Up
-- +goose StatementBegin
GRANT CONNECTION_ADMIN ON *.* TO 'mysql-tls-reloader';
-- +goose StatementEnd
-- +goose StatementBegin
GRANT ALL ON spiredemo.* TO 'spire-mysql-client';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
REVOKE CONNECTION_ADMIN ON *.* FROM 'mysql-tls-reloader';
-- +goose StatementEnd
-- +goose StatementBegin
REVOKE ALL ON spiredemo.* FROM 'spire-mysql-client';
-- +goose StatementEnd
