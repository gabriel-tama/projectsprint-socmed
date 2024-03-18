package friend

type FriendResponse struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int    `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}

type FriendListResponse []FriendResponse
