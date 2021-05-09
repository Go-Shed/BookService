package services

import (
	"shed/bookservice/api"
	"shed/bookservice/repos/dgraph/query"
)

type BookService struct {
	bookRepo query.BookRepo
}

func NewBookService() BookService {
	return BookService{bookRepo: query.NewBookRepo()}
}

func (p *BookService) GetBooks(userId string) ([]api.GetBooksResponse, error) {
	client := p.bookRepo
	user, err := client.GetBooks(userId)

	if err != nil {
		return []api.GetBooksResponse{}, err
	}

	var response []api.GetBooksResponse

	for _, book := range user.Books {
		item := api.GetBooksResponse{
			BookId:   book.Id,
			BookName: book.Name,
		}
		response = append(response, item)
	}

	return response, nil
}
