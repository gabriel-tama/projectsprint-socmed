package post

import (
	"context"

	"github.com/gabriel-tama/projectsprint-socmed/api/comment"
	"github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Create(ctx context.Context, post *Post) error
	GetAllPosts(ctx context.Context, req GetAllPostsPayload, userId int) (*GetAllPostsResponse, int, error)
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

func (d *dbRepository) GetAllPosts(ctx context.Context, req GetAllPostsPayload, userId int) (*GetAllPostsResponse, int, error) {
	var res GetAllPostsResponse
	stmt := `SELECT p.id, p.post_in_html, array_agg(t.tag) AS tags, p.created_at, c.content, c.created_at, u.id, u.name, COALESCE(u.imageUrl,''), u.friendsCount, pu.id, pu.name, COALESCE(pu.imageUrl,''),pu.friendsCount, pu.created_at, COUNT(*)
            FROM posts p
            JOIN users pu ON p.user_id = pu.id
            LEFT JOIN comments c ON p.id = c.post_id
            LEFT JOIN users u ON c.user_id = u.id
            LEFT JOIN post_tags t ON t.post_id = p.id
    `
	if req.Search != "" {
		stmt += `WHERE p.post_in_html LIKE %` + req.Search + `% `
	}
	if len(req.SearchTag) != 0 {
		if req.Search != "" {
			stmt += `AND ( `

		} else {
			stmt += `WHERE (`

		}
		for i := 0; i < len(req.SearchTag); i++ {
			stmt += `t.tag=` + req.SearchTag[i] + ` `
			if len(req.SearchTag)-1 != i {
				stmt += `OR `
			}
		}
		stmt += ` ) `
	}
	stmt += `
            GROUP BY p.id,c.id,u.id,pu.id LIMIT $1 OFFSET $2
    `
	rows, err := d.db.Pool.Query(ctx, stmt, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []comment.CommentResponse

	post_id_before := -1
	var total_count int
	for rows.Next() {
		var comment comment.CommentResponse
		var post PostResponse

		err := rows.Scan(&post.ID, &post.PostInHtml, &post.Tags, &post.CreatedAt,
			&comment.Content, &comment.CreatedAt, &comment.Creator.ID, &comment.Creator.Name, &comment.Creator.ImageUrl, &comment.Creator.FriendsCount,
			&post.Creator.ID, &post.Creator.Name, &post.Creator.ImageUrl, &post.Creator.FriendsCount, &post.Creator.CreatedAt,
			&total_count)
		if err != nil {
			return nil, 0, err
		}
		if post_id_before == -1 {
			res = append(res, post)
			post_id_before = post.ID
			comments = append(comments, comment)
			continue
		}
		if post_id_before != post.ID {
			res[len(res)-1].Comment = comments
			res = append(res, post)
			comments = nil
		}
		post_id_before = post.ID
		comments = append(comments, comment)
	}
	res[len(res)-1].Comment = comments

	return &res, total_count, err
}
