package post

import (
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/comment"
	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type PostResponse struct {
	ID         int                       `json:"postId"`
	PostInHtml string                    `json:"post_in_html"`
	Tags       []string                  `json:"tags"`
	CreatedAt  time.Time                 `json:"createdAt,omitempty"`
	Creator    user.UserResponse         `json:"creator,omitempty"`
	Comment    []comment.CommentResponse `json:"comments,omitempty"`
}

type GetAllPostsResponse []PostResponse
