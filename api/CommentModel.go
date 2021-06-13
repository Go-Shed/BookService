package api

type AddCommentRequest struct {
	Text   string `json:"text"`
	PostId string `json:"post_id"`
}

type GetCommentsRequest struct {
	PostId string `json:"post_id"`
}

type GetCommentsResponse struct {
	Comments []CommentItem `json:"comments"`
}

type CommentItem struct {
	Text      string `json:"text"`
	UserName  string `json:"user_name"`
	UserId    string `json:"user_id"`
	UserPhoto string `json:"user_photo"`
	CreatedAt string `json:"created_at"`
	CommentId string `json:"comment_id"`
}
