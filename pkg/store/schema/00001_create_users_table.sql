-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name varchar(25) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users;
-- +goose StatementEnd
