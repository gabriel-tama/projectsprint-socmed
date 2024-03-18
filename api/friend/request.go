package friend

type AddFriendPayload struct {
	UserId int `json:"userId" binding:"required"`
}

type DeleteFriendPayload struct {
	UserId int `json:"userId" binding:"required"`
}
