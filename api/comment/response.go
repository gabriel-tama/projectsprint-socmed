package comment

import (
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type CommentResponse struct {
	Content   string            `json:"comment"`
	Creator   user.UserResponse `json:"creator,omitempty"`
	CreatedAt time.Time         `json:"createdAt,omitempty"`
}
