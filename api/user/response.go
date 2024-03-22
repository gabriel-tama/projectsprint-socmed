package user

import "time"

type UserResponse struct {
	ID           string    `json:"userId,omitempty"`
	Name         string    `json:"name"`
	ImageUrl     string    `json:"imageUrl,omitempty"`
	AccessToken  string    `json:"accessToken,omitempty"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	FriendsCount int       `json:"friendCount,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}
