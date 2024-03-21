package friend

import "time"

type FriendResponse struct {
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
	ImageUrl    string    `json:"imageUrl"`
	FriendCount int       `json:"friendCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

type FriendListResponse []FriendResponse
