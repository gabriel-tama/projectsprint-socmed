package comment

import (
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/user"
)

type CommentResponse struct {
	PostId    int
	Content   string                `json:"comment"`
	Creator   user.RegisterResponse `json:"creator,omitempty"`
	CreatedAt time.Time             `json:"createdAt,omitempty"`
}
