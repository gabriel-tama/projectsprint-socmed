package post

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/comment"
	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type PostResponse struct {
	ID      string                    `json:"postId"`
	Post    SinglePost                `json:"post"`
	Creator user.UserResponse         `json:"creator,omitempty"`
	Comment []comment.CommentResponse `json:"comments,omitempty"`
}

type SinglePost struct {
	PostInHtml string    `json:"postInHtml"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}

type PostCreator struct {
	ID          NullString `json:"userId"`
	Name        NullString `json:"name"`
	ImageUrl    NullString `json:"imageUrl"`
	FriendCount NullInt64  `json:"friendCount"`
	CreatedAt   NullTime   `json:"createdAt"`
}

type PostComment struct {
	ID        NullString  `json:"id"`
	PostId    NullString  `json:"postId"`
	Content   NullString  `json:"comment"`
	CreatedAt NullTime    `json:"createdAt"`
	Creator   PostCreator `json:"creator"`
}

type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type NullTime struct {
	sql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

type GetAllPostsResponse []PostResponse
