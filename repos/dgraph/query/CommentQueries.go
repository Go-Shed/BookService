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

func (repo CommentRepo) GetComments(postId string) ([]model.Comment, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query GetComments {
			queryPost(filter: {id: "%s"}) {
			  comments(order: {desc: createdAt}) {
				id
				text
				createdAt
				user {
				  userId
				  userName
				  userPhoto
				}
			  }
			}
		  }
		  `, postId),
	}

	response, err := client.Do(query)

	if err != nil {
		return []model.Comment{}, err
	}

	var posts []model.Post
	mapstructure.Decode(response["queryPost"], &posts)

	if len(posts) == 0 {
		return []model.Comment{}, fmt.Errorf("comment not added")
	}

	return posts[0].Comments, nil
}

func (repo CommentRepo) GetTopCommentBulk(postIds []string) ([]model.Comment, error) {

	client := repo.client

	query := dgraph.Request{
		Query: `query GetComments($patch: [ID!]) {
			queryPost(filter: {id: $patch}) {
			  id
			  text
			  comments(order: {desc: createdAt}, first: 1) {
				id
				text
				createdAt
				user {
				  userId
				  userName
				  userPhoto
				}
			  }
			}
		  }`, Variables: dgraph.Variables{Patch: postIds},
	}

	response, err := client.Do(query)

	if err != nil {
		return []model.Comment{}, err
	}

	var posts []model.Post
	mapstructure.Decode(response["queryPost"], &posts)

	var result []model.Comment

	for _, post := range posts {
		if len(post.Comments) > 0 {
			result = append(result, post.Comments[0])
		} else {
			result = append(result, model.Comment{})
		}
	}

	return result, nil
}
