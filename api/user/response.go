package user

type UserResponse struct {
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
}
