package api

type FollowUserRequest struct {
	UserId string `json:"user_id"`
}

type SearchUserRequest struct {
	UserName string `json:"user_name,omitempty"`
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
