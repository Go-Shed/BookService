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

	var user model.User
	mapstructure.Decode(response["queryUser"], &user)
	return user, nil
}
