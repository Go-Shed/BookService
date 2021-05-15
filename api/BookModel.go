package api

type GetBooksRequest struct {
	ProfileUserId string `json:"profile_user_id"`
	IsSelf        bool   `json:"is_self"`
}

type GetBooksResponse struct {
	BookId   string `json:"book_id"`
	BookName string `json:"book_name"`
}
