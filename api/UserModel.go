package api

type FollowUserRequest struct {
	UserId string `json:"user_id"`
}

type SearchUserRequest struct {
	Search string `json:"search,omitempty"`
}

type SearchUserResponse struct {
	SearchResults []SearchResult `json:"search_results"`
}

type SearchResult struct {
	UserPhoto        string `json:"user_photo"`
	ShowFollowButton bool   `json:"show_follow_button"`
	IsFollowed       bool   `json:"is_followed"`
	UserName         string `json:"user_name"`
	UserId           string `json:"user_id"`
}

type FetchProfileResponse struct {
	IsFollowing bool   `json:"is_following"`
	UserPhoto   string `json:"user_photo,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Followers   int64  `json:"followers"`
	Following   int64  `json:"following"`
}

type FetchProfileRequest struct {
	IsSelf        bool   `json:"is_self"`
	ProfileUserId string `json:"profile_user_id,omitempty"`
}

type GetFollowersResponse struct {
	Follows        []FollowItem `json:"follows"`
	TotalFollowers int          `json:"total_followers"`
}

type GetFollowingResponse struct {
	Follows        []FollowItem `json:"follows"`
	TotalFollowing int          `json:"total_following"`
}

type FollowItem struct {
	UserPhoto  string `json:"user_photo"`
	UserName   string `json:"user_name"`
	UserId     string `json:"user_id"`
	IsFollowed bool   `json:"is_followed"`
}
