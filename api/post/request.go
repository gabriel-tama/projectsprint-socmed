package post

type CreatePostPayload struct {
	PostInHTML string   `json:"postInHtml" binding:"required,min=0,max=500"`
	Tags       []string `json:"tags" binding:"required,dive,required"`
}

func (p CreatePostPayload) Validate() error {
	// Check if PostInHTML is not empty, has minimum length 2, and maximum length 500
	// if len(strings.TrimSpace(p.PostInHTML)) < 2 || len(p.PostInHTML) > 500 {
	// 	return fmt.Errorf("postInHtml must be between 2 and 500 characters")
	// }

	// Check if Tags is not empty
	// if len(p.Tags) == 0 {
	// 	return fmt.Errorf("tags must not be empty")
	// }

	// You can add more validations for the Tags if needed

	return nil
}

type GetAllPostsPayload struct {
	Limit     int      `form:"limit,default=5" binding:"min=0"`
	Offset    int      `form:"offset,default=0" binding:"min=0"`
	Search    string   `form:"search"`
	SearchTag []string `form:"searchTag"`
}
