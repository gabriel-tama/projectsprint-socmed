package comment

import (
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	User      user.RegisterResponse
	CreatedAt time.Time
}
