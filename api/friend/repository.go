package friend

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	AddFriend(ctx context.Context, userId int, requestedId int) error
	DeleteFriend(ctx context.Context, userId int, requestedId int) error
	GetAllFriends(ctx context.Context, userId int, req GetAllFriendsPayload) (*FriendListResponse, error, int)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) AddFriend(ctx context.Context, userId int, requestedId int) error {
	// var exists bool
	var pgErr *pgconn.PgError

	// row := d.db.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM friends WHERE user_id=$1 AND friend_id=$2)", requestedId, userId)
	// err := row.Scan(&exists)
	// if err != nil {
	// 	return err
	// }

	// if exists {
	// 	return ErrAlreadyFriends
	// }

	err := d.db.StartTx(ctx, func(tx pgx.Tx) error {
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
		_, err = tx.Exec(ctx, "INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)", requestedId, userId)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "UPDATE users SET friendsCount=friendsCount+1 WHERE id=$1 OR id=$2", userId, requestedId)
		return err
	})
	return err
}

func (d *dbRepository) DeleteFriend(ctx context.Context, userId int, requestedId int) error {
	// row, err := d.db.Pool.Exec(ctx, "DELETE FROM friends WHERE (user_id=$1 AND friend_id=$2) OR (user_id=$2 AND friend_id=$1)", userId, requestedId)
	// if err != nil {
	// 	return err
	// }

	// if row.RowsAffected() == 0 {
	// 	return ErrNotFriends
	// }

	err := d.db.StartTx(ctx, func(tx pgx.Tx) error {
		row, err := tx.Exec(ctx, "DELETE FROM friends WHERE (user_id=$1 AND friend_id=$2) OR (user_id=$2 AND friend_id=$1)", userId, requestedId)
		if err != nil {
			return err
		}
		if row.RowsAffected() == 0 {
			return ErrNotFriends
		}
		_, err = tx.Exec(ctx, "UPDATE users SET friendsCount=friendsCount-1 WHERE id=$1 OR id=$2", requestedId, userId)
		return err
	})

	return err
}

func (d *dbRepository) GetAllFriends(ctx context.Context, userId int, req GetAllFriendsPayload) (*FriendListResponse, error, int) {
	var friendsList FriendListResponse
	total := 0
	// stmt := `SELECT u.id,u.name,COALESCE(u.imageUrl,''),u.friendsCount,u.created_at,COUNT(*) AS total_count FROM users AS u `

	// if req.OnlyFriend == true {
	// 	stmt += `LEFT JOIN friends f ON f.friend_id=u.id
	// 			WHERE u.id!=$1 `
	// } else {
	// 	stmt += "WHERE u.id!=$1 "
	// }

	// if req.Search != "" {
	// 	stmt += `AND u.name LIKE %` + req.Search + `% `
	// }
	// stmt += `GROUP BY u.id `
	stmt := `SELECT u.id, u.name, COALESCE(u.imageUrl,''),u.friendsCount,u.created_at,COUNT(*) OVER() FROM users u `
	if req.OnlyFriend == true {
		stmt += `WHERE u.id IN (SELECT friend_id FROM friends WHERE user_id=` + strconv.Itoa(userId) + `) `
		if req.Search != "" {
			stmt += `AND u.name LIKE '%` + req.Search + `%' `
		}
	} else {
		stmt += `WHERE u.id!=` + strconv.Itoa(userId)
		if req.Search != "" {
			stmt += `AND u.name LIKE '%` + req.Search + `%' `
		}
	}
	stmt += `GROUP BY u.id `
	if req.SortBy == "createdAt" {
		stmt += `ORDER BY created_at `

	} else {
		stmt += `ORDER BY friendsCount `

	}

	if req.OrderBy == "asc" {
		stmt += `ASC `
	} else {
		stmt += `DESC `
	}
	stmt += `LIMIT $1 OFFSET $2 `
	fmt.Println(stmt)
	rows, err := d.db.Pool.Query(ctx, stmt, req.Limit, req.Offset)
	if err != nil {
		return nil, err, 0
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println(rows)
		var friend FriendResponse
		var tot int
		err := rows.Scan(&friend.UserId, &friend.Name, &friend.ImageUrl, &friend.FriendCount, &friend.CreatedAt, &tot)
		if err != nil {
			return nil, err, 0
		}
		friendsList = append(friendsList, friend)
		total += tot
	}
	if err := rows.Err(); err != nil {
		return nil, err, 0
	}

	return &friendsList, nil, total
}
