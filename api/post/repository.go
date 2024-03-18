package post

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, post *Post) error
}

type dbRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool, BCRYPT_SALT int) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, post *Post) error {

	return nil
}
