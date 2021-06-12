package api

type AddCommentRequest struct {
	Text   string `json:"text"`
	PostId string `json:"post_id"`
}
