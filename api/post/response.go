package post

import (
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type PostResponse struct {
	PostInHtml string     `json:"post_in_html"`
	Tags       []string   `json:"tags"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	Creator    *user.User `json:"creator,omitempty"`
}
