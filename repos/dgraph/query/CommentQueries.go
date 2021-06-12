package query

import (
	"fmt"
	"shed/bookservice/repos/dgraph"
	"shed/bookservice/repos/dgraph/model"

	"github.com/mitchellh/mapstructure"
)

type CommentRepo struct {
	client dgraph.Dgraph
}

func NewCommentRepo() CommentRepo {
	return CommentRepo{client: dgraph.Dgraph{}}
}

func (repo CommentRepo) AddComment(comment model.Comment) error {

	client := repo.client
	query := dgraph.Request{
		Query: `mutation addComment($patch: [AddCommentInput!]!) {
			addComment(input: $patch) {
			  comment {
				id
			  }
			}
		  }`,
		Variables: dgraph.Variables{Patch: comment}, Operation: "addComment"}

	response, err := client.Do(query)

	if err != nil {
		return err
	}

	var comments []model.Comment
	data := response["addComment"].(map[string]interface{})
	mapstructure.Decode(data["comment"], &comments)

	if len(comments) == 0 {
		return fmt.Errorf("comment not added")
	}

	return nil
}
