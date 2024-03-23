package friend

type AddFriendPayload struct {
	UserId string `json:"userId" binding:"required"`

}

type DeleteFriendPayload struct {
	UserId  string `json:"userId" binding:"required"`
	UserInt int
}

type GetAllFriendsPayload struct {
	Limit      int    `form:"limit,default=5" binding:"min=0"`
	Offset     int    `form:"offset,default=0" binding:"min=0"`
	SortBy     string `form:"sortBy,default=createdAt" binding:"oneof=friendCount createdAt"`
	OrderBy    string `form:"orderBy,default=desc" binding:"oneof=desc asc"`
	OnlyFriend bool   `form:"onlyFriend,default=false" binding:"boolean"`
	Search     string `form:"search"`
}
