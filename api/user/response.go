package user

import "time"

type RegisterResponse struct {
	ID           string    `json:"userId,omitempty"`
	Name         string    `json:"name"`
	ImageUrl     string    `json:"imageUrl,omitempty"`
	AccessToken  string    `json:"accessToken,omitempty"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	FriendsCount string    `json:"friendCount,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}

type LoginResponse struct {
	ID           string    `json:"userId,omitempty"`
	Name         string    `json:"name"`
	ImageUrl     string    `json:"imageUrl,omitempty"`
	AccessToken  string    `json:"accessToken,omitempty"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	FriendsCount string    `json:"friendCount,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}
