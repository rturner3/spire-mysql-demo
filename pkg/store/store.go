package store

import (
	"context"
	"database/sql"
	"log"
)

const (
	listUsersQuery  = "SELECT * FROM Users"
	createUserQuery = "INSERT INTO Users (name) VALUES ( ? );"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

type User struct {
	ID   int
	Name string
}

func (s *Store) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, listUsersQuery)
	if err != nil {
		log.Printf("Failed to run list users query: %v", err)
		return nil, err
	}

	users := make([]User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name)
		users = append(users, user)
	}

	return users, nil
}

func (s *Store) CreateUser(ctx context.Context, user User) error {
	if _, err := s.db.ExecContext(ctx, createUserQuery, user.Name); err != nil {
		log.Printf("Failed to run list users query: %v", err)
		return err
	}
	return nil
}
