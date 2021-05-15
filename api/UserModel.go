package api

type FollowUserRequest struct {
	UserId string `json:"user_id,omitempty"`
}

type SearchUserRequest struct {
	UserName string `json:"user_name,omitempty"`
}

type FetchProfileResponse struct {
	UserPhoto   string `json:"user_photo,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Followers   int64  `json:"followers,omitempty"`
	Following   int64  `json:"following,omitempty"`
	IsFollowing bool   `json:"is_following,omitempty"`
}

type FetchProfileRequest struct {
	IsSelf        bool   `json:"is_self,omitempty"`
	ProfileUserId string `json:"profile_user_id,omitempty"`
}
