package query

import (
	"fmt"
	"shed/bookservice/repos/dgraph"
	"shed/bookservice/repos/dgraph/model"

	"github.com/mitchellh/mapstructure"
)

type BookRepo struct {
	client dgraph.Dgraph
}

func NewBookRepo() BookRepo {
	return BookRepo{client: dgraph.Dgraph{}}
}

func (repo BookRepo) GetBooks(userId string) (model.User, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query{
			queryUser(filter: {userId: {eq: "%s"}}){
			  books{
				id
				name
			  }
			}
		  }`, userId),
	}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, err
	}

	var user []model.User
	mapstructure.Decode(response["queryUser"], &user)

	if len(user) == 0 {
		return model.User{}, nil
	}

	return user[0], nil
}

func (repo BookRepo) CreateOrGetBook(bookId, bookName string) (string, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query{
			queryBook(filter: {name: {eq: "%s"}, or: {id: "%s"}}) {
			  name
			  id
			}
		  }`, bookName, bookId),
	}

	response, err := client.Do(query)

	if err != nil {
		return "", err
	}

	var books []model.Book
	mapstructure.Decode(response["queryBook"], &books)

	if len(books) == 0 {
		return repo.CreateBook(bookName)
	}

	return books[0].Id, nil
}

func (repo BookRepo) CreateBook(bookName string) (string, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			addBook(input: {name: "%s"}){
			  book{
				id
			  }
			}
		  }
		  `, bookName),
	}

	response, err := client.Do(query)

	if err != nil {
		return "", err
	}

	var book []model.Book
	data := response["addBook"].(map[string]interface{})
	mapstructure.Decode(data["book"], &book)

	return book[0].Id, nil
}
