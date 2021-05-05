package api

type GetPostsRequest struct {
	UserId     string `json:"user_id"`
	ScreenName string `json:"screen_name"`
}

type GetPostsResponse struct {
	UserId          string `json:"user_id"`
	UserPhoto       string `json:"user_photo"`
	Text            string `json:"text"`
	PostColor       string `json:"post_color"`
	ShowFollowBtn   bool   `json:"show_follow_button"`
	ShowEditBtn     bool   `json:"show_edit_button"`
	ShowLikeSection bool   `json:"show_like_section"`
	IsFollowed      bool   `json:"is_followed"`
	IsLiked         bool   `json:"is_liked"`
	LikeCount       string `json:"like_count"`
	UserName        string `json:"user_name"`
	PostId          string `json:"post_id"`
	CreatedAt       string `json:"created_at"`
}

type AddPostRequest struct {
	Text      string `json:"text"`
	PostColor string `json:"post_color"`
}

type UpdatePostRequest struct {
	Text   string `json:"text"`
	PostId string `json:"post_id"`
}
