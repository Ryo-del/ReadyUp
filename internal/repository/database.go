package repository

import (
	"context"

	"ReadyUp/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, email, username, passwordHash string) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, email, username, passwordHash string) error {
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1,$2,$3)
	`
	_, err := r.db.Exec(ctx, query, email, username, passwordHash)
	if err != nil {
		return err
	}
	return nil
}
