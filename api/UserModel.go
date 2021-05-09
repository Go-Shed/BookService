package api

type FollowUserRequest struct {
	UserId string `json:"user_id,omitempty"`
}

type SearchUserRequest struct {
	UserName string `json:"user_name,omitempty"`
}
