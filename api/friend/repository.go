package friend

import (
	"context"
	"errors"

	"github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	AddFriend(ctx context.Context, userId int, requestedId int) error
	DeleteFriend(ctx context.Context, userId int, requestedId int) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) AddFriend(ctx context.Context, userId int, requestedId int) error {
	var exists bool
	var pgErr *pgconn.PgError

	row := d.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM friends WHERE user_id=$1 AND friend_id=$2)", requestedId, userId)
	err := row.Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyFriends
	}

	err = d.db.StartTx(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)", userId, requestedId)
		if err != nil {
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case "23505":
					return ErrAlreadyFriends
				case "23503":
					return ErrInvalidUser
				default:
					return err
				}
			}
			return err
		}
		return err
	})
	return err
}

func (d *dbRepository) DeleteFriend(ctx context.Context, userId int, requestedId int) error {
	row, err := d.db.Pool.Exec(ctx, "DELETE FROM friends WHERE (user_id=$1 AND friend_id=$2) OR (user_id=$2 AND friend_id=$1)", userId, requestedId)
	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		return ErrNotFriends
	}

	return nil
}
