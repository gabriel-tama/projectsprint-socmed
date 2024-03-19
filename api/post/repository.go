package post

import (
	"context"

	"github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Create(ctx context.Context, post *Post) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, post *Post) error {
	err := d.db.StartTx(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "INSERT INTO posts (user_id, post_in_html) VALUES ($1, $2) RETURNING id", post.UserId, post.PostInHtml).Scan(&post.ID)
		if err != nil {
			return err
		}
		stmt := `INSERT INTO post_tags (post_id, tag) VALUES ($1, $2)`
		// Insert each tag
		for _, tag := range post.Tags {
			_, err := tx.Exec(ctx, stmt, post.ID, tag)
			if err != nil {
				return err
			}
			post.Tags = append(post.Tags, tag)
		}

		return err
	})
	if err != nil {
		return err
	}
	return nil
}
