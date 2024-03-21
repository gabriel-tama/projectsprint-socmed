package comment

import (
	"context"
	"errors"

	"github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Create(ctx context.Context, comment *CreateCommentPayload, user_id int) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, comment *CreateCommentPayload, user_id int) error {

	var exists bool
	var postCreatorId int

	row := d.db.Pool.QueryRow(ctx, "SELECT user_id FROM posts WHERE id=$1", comment.PostID)
	err := row.Scan(&postCreatorId)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidPost
	}

	if err != nil {
		return err
	}

	row = d.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM friends WHERE user_id=$1 AND friend_id=$2)", postCreatorId, user_id)
	err = row.Scan(&exists)

	if err != nil {
		return err
	}
	if !exists && (postCreatorId != user_id) {
		return ErrNotFriends
	}

	err = d.db.StartTx(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO comments (user_id, post_id, content) VALUES ($1,$2,$3)", user_id, comment.PostID, comment.Content)
		if err != nil {
			return err
		}
		return err
	})
	return err
}
