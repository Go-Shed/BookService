package api

type GetBooksRequest struct {
	ProfileUserId string `json:"profile_user_id"`
	IsSelf        bool   `json:"is_self"`
	ScreenName    string `json:"screen_name"`
}

type GetBooksResponse struct {
	Books []BookResponse `json:"books"`
}

type BookResponse struct {
	BookId   string `json:"book_id"`
	BookName string `json:"book_name"`
}
