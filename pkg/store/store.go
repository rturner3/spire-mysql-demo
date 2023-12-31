package store

import (
	"context"
	"database/sql"
	"log"
	"sync"
)

const (
	listUsersQuery  = "SELECT * FROM Users"
	createUserQuery = "INSERT INTO Users (name) VALUES ( ? );"
)

type Store struct {
	mu sync.RWMutex
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	rows, err := s.db.QueryContext(ctx, listUsersQuery)
	if err != nil {
		log.Printf("Failed to run list users query: %v", err)
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name)
		users = append(users, user)
	}

	return users, nil
}

func (s *Store) CreateUser(ctx context.Context, user User) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, err := s.db.ExecContext(ctx, createUserQuery, user.Name); err != nil {
		log.Printf("Failed to run list users query: %v", err)
		return err
	}
	return nil
}

func (s *Store) UpdateDB(db *sql.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.db.Close(); err != nil {
		log.Printf("Failed to close existing DB connection: %v", err)
	}
	s.db = db
}
