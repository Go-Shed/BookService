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
	Text      string `json:"text,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	UserPhoto string `json:"user_photo,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	CommentId string `json:"comment_id,omitempty"`
}
