package post

import (
	"context"
	"fmt"

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

// func (d *dbRepository) GetAllPosts(ctx context.Context, req GetAllPostsPayload, userId int) (*GetAllPostsResponse, int, error) {
// 	var res GetAllPostsResponse
// 	stmt := `SELECT p.id, p.post_in_html, array_agg(t.tag) AS tags, p.created_at, c.content, c.created_at, u.id, u.name, COALESCE(u.imageUrl,''), u.friendsCount, pu.id, pu.name, COALESCE(pu.imageUrl,''),pu.friendsCount, pu.created_at, COUNT(*)
//             FROM posts p
//             JOIN users pu ON p.user_id = pu.id
//             LEFT JOIN comments c ON p.id = c.post_id
//             LEFT JOIN users u ON c.user_id = u.id
//             LEFT JOIN post_tags t ON t.post_id = p.id
//     `
// 	if req.Search != "" {
// 		stmt += `WHERE p.post_in_html LIKE %` + req.Search + `% `
// 	}
// 	if len(req.SearchTag) != 0 {
// 		if req.Search != "" {
// 			stmt += `AND ( `

// 		} else {
// 			stmt += `WHERE (`

// 		}
// 		for i := 0; i < len(req.SearchTag); i++ {
// 			stmt += `t.tag=` + req.SearchTag[i] + ` `
// 			if len(req.SearchTag)-1 != i {
// 				stmt += `OR `
// 			}
// 		}
// 		stmt += ` ) `
// 	}
// 	stmt += `
//             GROUP BY p.id,c.id,u.id,pu.id LIMIT $1 OFFSET $2
//     `
// 	fmt.Println(stmt)
// 	rows, err := d.db.Pool.Query(ctx, stmt, req.Limit, req.Offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer rows.Close()

// 	var comments []comment.CommentResponse
// 	// fmt.Println(rows)
// 	post_id_before := -1
// 	var total_count int
// 	for rows.Next() {
// 		var comment comment.CommentResponse
// 		var post PostResponse
// 		// var commentRes comment.CommentResponse

// 		err = rows.Scan(&post.ID, &post.PostInHtml, &post.Tags, &post.CreatedAt,
// 			&comment.Content, &comment.CreatedAt, &comment.Creator.ID, &comment.Creator.Name, &comment.Creator.ImageUrl, &comment.Creator.FriendsCount,
// 			&post.Creator.ID, &post.Creator.Name, &post.Creator.ImageUrl, &post.Creator.FriendsCount, &post.Creator.CreatedAt,
// 			&total_count)
// 		if err != nil {
// 			// log.Fatal(err)
// 			return nil, 0, err
// 		}
// 		// articleJSON, err := json.Marshal(&comment)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return nil, 0, err
// 		// }
// 		// log.Printf("json marshal := %s\n\n", articleJSON)
// 		// fmt.Println(post.ID)

// 		if post_id_before == -1 {
// 			res = append(res, post)
// 			post_id_before = post.ID
// 			comments = append(comments, comment)
// 			continue
// 		}

// 		if post_id_before != post.ID {
// 			res[len(res)-1].Comment = comments
// 			res = append(res, post)
// 			comments = nil
// 		}

// 		post_id_before = post.ID
// 		// if comment.ID==-1{
// 		//     continue
// 		// }
// 		comments = append(comments, comment)
// 	}
// 	res[len(res)-1].Comment = comments
// 	// if rows.Err() != nil {
// 	// 	fmt.Printf("rows error: %v", rows.Err())
// 	// 	return nil, total_count, err
// 	// }

// 	return &res, total_count, err
// }

func (d *dbRepository) GetAllPosts(ctx context.Context, req GetAllPostsPayload, userId int) (*GetAllPostsResponse, int, error) {
	var res GetAllPostsResponse
	stmt := `WITH UserPosts AS (
                SELECT p.id AS post_id, p.user_id, p.post_in_html,p.created_at AS post_created_at, 
				u.name, COALESCE(u.imageUrl,''), u.created_at, ARRAY_AGG(pt.tag) AS tags,u.friendsCount,
				COALESCE(c.post_id, -1), COALESCE(c.content,''),COALESCE(c.created_at,'2024-03-24 04:51:25.503309'),
				COALESCE(cu.id,-1), COALESCE(cu.name,''), COALESCE(cu.friendsCount,0), COALESCE(cu.imageUrl,'')
                FROM posts p
                LEFT JOIN post_tags pt ON p.id = pt.post_id
                LEFT JOIN users u ON p.user_id = u.id
				LEFT JOIN comments c ON p.id = c.post_id
				LEFT JOIN users cu ON c.user_id = cu.id
                WHERE p.user_id = $1
                GROUP BY p.id, p.user_id, u.id, c.id, cu.id
                
                UNION ALL
                
                SELECT p.id AS post_id, p.user_id, p.post_in_html,p.created_at AS post_created_at,
				u.name, COALESCE(u.imageUrl,''), u.created_at, ARRAY_AGG(pt.tag) AS tags,u.friendsCount,
				COALESCE(c.post_id, -1), COALESCE(c.content,''),COALESCE(c.created_at,'2024-03-24 04:51:25.503309'),
				COALESCE(cu.id,-1), COALESCE(cu.name,''), COALESCE(cu.friendsCount,0), COALESCE(cu.imageUrl,'')                
				FROM posts p
                LEFT JOIN post_tags pt ON p.id = pt.post_id
                LEFT JOIN users u ON p.user_id = u.id
				LEFT JOIN comments c ON p.id = c.post_id
				LEFT JOIN users cu ON c.user_id = cu.id
                WHERE EXISTS (
                    SELECT 1
                    FROM friends f
                    WHERE f.user_id = $1 AND f.friend_id = p.user_id
                )
                GROUP BY p.id, p.user_id, u.id, c.id, cu.id
            )
            SELECT *,
                (SELECT COUNT(DISTINCT post_id) FROM UserPosts
				
				
             `
	filter := ``
	if req.Search != "" {
		filter += "WHERE post_in_html LIKE '%" + req.Search + "%' "
	}
	if len(req.SearchTag) != 0 {
		if req.Search != "" {
			filter += "AND ( "
		} else {
			filter += "WHERE "
		}
		for i := 0; i < len(req.SearchTag); i++ {
			filter += `'` + req.SearchTag[i] + `' = ANY(tags) `
			if len(req.SearchTag)-1 != i {
				filter += `OR `
			}
		}
		if req.Search != "" {
			filter += ") "
		}
	}
	stmt += filter

	stmt += `) AS total_post_count
            FROM UserPosts `
	stmt += filter

	stmt += "ORDER BY post_created_at DESC LIMIT $2 OFFSET $3 "
	fmt.Println(stmt)
	rows, err := d.db.Pool.Query(ctx, stmt, userId, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(userId)
	defer rows.Close()
	var total int
	var postIds []string
	var co_postId []string
	var comments []comment.CommentResponse
	for rows.Next() {
		var post PostResponse
		var comment comment.CommentResponse
		var c_postId string
		err = rows.Scan(&post.ID, &post.Creator.ID, &post.Post.PostInHtml, &post.Post.CreatedAt, &post.Creator.Name, &post.Creator.ImageUrl, &post.Creator.CreatedAt, &post.Post.Tags, &post.Creator.FriendsCount,
			&c_postId, &comment.Content, &comment.CreatedAt,
			&comment.Creator.ID, &comment.Creator.Name, &comment.Creator.FriendsCount, &comment.Creator.ImageUrl,
			&total)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, post)
		postIds = append(postIds, post.ID)
		if c_postId != "-1" {
			comments = append(comments, comment)
			co_postId = append(co_postId, c_postId)
		}
	}
	var new_res GetAllPostsResponse
	idx_post := 0
	for i := 0; i < len(res); i++ {
		if res[i].ID == res[idx_post].ID {
			new_res = append(new_res, res[i])
			idx_post++
		}
	}
	idx_post = 0
	for i := 0; i < len(co_postId); i++ {
		if co_postId[i] == new_res[idx_post].ID {
			new_res[idx_post].Comment = append(new_res[idx_post].Comment, comments[i])
		} else {
			idx_post++
			new_res[idx_post].Comment = append(new_res[idx_post].Comment, comments[i])
		}
	}

	// stmt = `SELECT c.content,c.created_at,c.user_id,u.name,COALESCE(u.imageUrl,''),u.friendsCount
	//         FROM comments c
	//         LEFT JOIN users u on c.user_id = u.id
	//         WHERE c.post_id=$1
	//         GROUP BY c.post_id,c.content,c.created_at,c.user_id,u.name,u.imageUrl,u.friendsCount
	//         ORDER BY c.created_at DESC`
	// fmt.Println(postIds)

	// for i := 0; i < len(postIds); i++ {
	// 	rows, err := d.db.Pool.Query(ctx, stmt, postIds[i])
	// 	if err != nil {
	// 		return nil, 0, err
	// 	}
	// 	defer rows.Close()
	// 	for rows.Next() {
	// 		var comment comment.CommentResponse
	// 		err = rows.Scan(&comment.Content, &comment.CreatedAt, &comment.Creator.ID, &comment.Creator.Name, &comment.Creator.ImageUrl, &comment.Creator.FriendsCount)
	// 		if err != nil {
	// 			return nil, 0, err
	// 		}
	// 		// fmt.Println(comment)
	// 		res[i].Comment = append(res[i].Comment, comment)
	// 	}
	// }
	return &new_res, total, err
}
