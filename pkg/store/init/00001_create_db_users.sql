-- +goose Up
-- +goose StatementBegin
CREATE USER 'mysql-tls-reloader' REQUIRE SUBJECT '/C=US/O=SPIRE/CN=tls-reloader';
-- +goose StatementEnd
-- +goose StatementBegin
CREATE USER 'spire-mysql-client' REQUIRE SUBJECT '/C=US/O=SPIRE/CN=spire-mysql-client';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP USER 'mysql-tls-reloader';
-- +goose StatementEnd
-- +goose StatementBegin
DROP USER 'spire-mysql-client';
-- +goose StatementEnd
