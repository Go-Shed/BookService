package api

type GetBooksRequest struct {
	UserId string `json:"user_id"`
}

type GetBooksResponse struct {
	BookId   string `json:"book_id"`
	BookName string `json:"book_name"`
}
