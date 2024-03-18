package post

import (
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type Post struct {
	ID         uint64
	PostInHtml string
	Tags       []string
	CreatedAt  time.Time
	Creator    *user.User
}
