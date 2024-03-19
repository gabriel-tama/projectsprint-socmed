package comment

type CreateCommentPayload struct {
	PostID  string `json:"postId" binding:"required,alphanum"`
	Content string `json:"comment" binding:"required,min=2,max=500"`
}
