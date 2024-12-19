package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/wahyuoi/sbc/internal/model"
)

var _ UserRepository = (*sqlUserRepository)(nil)

type sqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) UserRepository {
	return &sqlUserRepository{db: db}
}

func (r *sqlUserRepository) Create(ctx context.Context, user *model.User) error {
	executor := getExecutor(ctx, r.db)

	query := `INSERT INTO users (email, password, created_at, updated_at) 
			  VALUES (?, ?, NOW(), NOW())`

	result, err := executor.ExecContext(ctx, query, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *sqlUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	executor := getExecutor(ctx, r.db)

	user := &model.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?`

	err := executor.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *sqlUserRepository) GetById(ctx context.Context, id int) (*model.User, error) {
	executor := getExecutor(ctx, r.db)

	user := &model.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = ?`

	err := executor.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
