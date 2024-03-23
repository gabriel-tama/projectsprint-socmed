package user

import "time"

type UserResponse struct {
	ID           string    `json:"userId,omitempty"`
	Name         string    `json:"name"`
	ImageUrl     string    `json:"imageUrl,omitempty"`
	AccessToken  string    `json:"accessToken,omitempty"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	FriendsCount string    `json:"friendCount,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}
