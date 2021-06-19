package api

type GetPostsRequest struct {
	ScreenName string `json:"screen_name"`
	ForUserID  string `json:"profile_user_id"`
	IsSelf     bool   `json:"is_self"`
}

type GetPostsResponse struct {
	Posts []GetPostResponse `json:"posts"`
}

type GetPostResponse struct {
	UserId          string       `json:"user_id"`
	UserPhoto       string       `json:"user_photo"`
	Text            string       `json:"text"`
	ShowFollowBtn   bool         `json:"show_follow_button"`
	ShowEditBtn     bool         `json:"show_edit_button"`
	ShowLikeSection bool         `json:"show_like_section"`
	IsFollowed      bool         `json:"is_followed"`
	IsLiked         bool         `json:"is_liked"`
	LikeCount       string       `json:"like_count"`
	UserName        string       `json:"user_name"`
	PostId          string       `json:"post_id"`
	CreatedAt       string       `json:"created_at"`
	Book            BookResponse `json:"book"`
	TopComment      CommentItem  `json:"top_comment"`
}

type AddPostRequest struct {
	Post      string `json:"post"`
	BookTitle string `json:"book_title"`
	BookId    string `json:"book_id"`
}

type UpdatePostRequest struct {
	Post      string `json:"post"`
	PostId    string `json:"post_id"`
	BookTitle string `json:"book_title"`
	BookId    string `json:"book_id"`
}

type PostIdRequest struct {
	PostId string `json:"post_id"`
}

type GetLikesResponse struct {
	Likes      []LikeItem `json:"likes"`
	TotalLikes int        `json:"total_likes"`
}

type LikeItem struct {
	UserPhoto     string `json:"user_photo"`
	ShowFollowBtn bool   `json:"show_follow_btn"`
	IsFollowed    bool   `json:"is_followed"`
	UserName      string `json:"user_name"`
	UserId        string `json:"user_id"`
}
