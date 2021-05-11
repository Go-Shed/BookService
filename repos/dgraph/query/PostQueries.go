package query

import (
	"fmt"
	"shed/bookservice/repos/dgraph"
	"shed/bookservice/repos/dgraph/model"

	"github.com/mitchellh/mapstructure"
)

type PostRepo struct {
	client dgraph.Dgraph
}

func NewPostRepo() PostRepo {
	return PostRepo{client: dgraph.Dgraph{}}
}

func (repo PostRepo) GetUserPosts(userId string) model.User {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query {
			getUser(username: "%s") {
			  posts{
				  id
				  title
				  text
				  book{
					name
					id
				  }
			  }
			}
		  }`, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user
}

func (repo PostRepo) GetUserFeedHomeScreen(userId string) (model.User, error) {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query MyQuery {
			getUser(userId: "%s") {
			  following {
				userId
				userName
				userPhoto
				posts {
				  id
				  text
				  createdAt
				  likes(filter: {userId: {eq: "%s"}}) {
					userId
				  }
				  likesAggregate{
					  count
				  }
				  book{
					name
					id
				}
				}
			  }
			}
		}`, userId, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, nil
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user, nil
}

func (repo PostRepo) GetUserHomeProfileScreen(userId, forUserId string) (model.User, error) {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`query MyQuery {
			getUser(userId: "%s") {
				userName
				userPhoto
    		posts {
				  id
				  text
				  createdAt
				  likes(filter: {userId: {eq: "%s"}}) {
					userId
				  }
				  followers(filter: {userId: {eq: "%s"}}) {
					userId
				  }
				  likesAggregate{
					  count
				  }
				  book{
					name
					id
				}
        }
			}
}`, forUserId, userId, userId)}

	response, err := client.Do(query)

	if err != nil {
		return model.User{}, nil
	}

	var user model.User
	mapstructure.Decode(response["getUser"], &user)
	return user, nil
}

func (repo PostRepo) GetExporeScreen(userId string) ([]model.Post, error) {
	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`{
			queryPost {
			  id
			  author {
				userId
				userPhoto
				userName
				followers(filter: {userId: {eq: "%s"}}){
				  userId
				}
			  }
			  text
			  likesAggregate {
				count
			  }
			  likes(filter: {userId: {eq: "%s"}}) {
				userId
			  }
			  book{
				name
				id
				}
			}
		  }
		  `, userId, userId),
	}

	response, err := client.Do(query)

	if err != nil {
		return []model.Post{}, err
	}

	var posts []model.Post
	mapstructure.Decode(response["queryPost"], &posts)
	return posts, nil
}

func (repo PostRepo) CreatePost(user model.User) error {

	type Userid struct {
		Eq string `json:"eq"`
	}

	type Filter struct {
		UserId Userid `json:"userId"`
	}

	type Set struct {
		Books model.Book `json:"books"`
		Posts model.Post `json:"posts"`
	}

	type Patch struct {
		Filter Filter `json:"filter"`
		Set    Set    `json:"set"`
	}

	patch := Patch{
		Filter: Filter{
			UserId: Userid{
				Eq: user.UserId,
			},
		},

		Set: Set{
			Books: user.Books[0],
			Posts: user.Posts[0],
		},
	}

	client := repo.client
	query := dgraph.Request{
		Query: `mutation addPostForUser($patch: UpdateUserInput!) {
			updateUser(input: $patch) {
    		user{
       			posts{
          		id
        	}
     	 }
		}
	}`, Variables: dgraph.Variables{Patch: patch}, Operation: "updateUser"}

	_, err := client.Do(query)

	if err != nil {
		return err
	}

	return nil
}

func (repo PostRepo) UpdatePost(post model.Post) (model.Post, error) {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation{
			updatePost(input: {filter: {id: "%s"}, set: {text: "%s", updatedAt: "%s"}}){
			  post{
				id
			  }
			}
		  }`, post.Id, post.Text, post.UpdatedAt), Operation: "updateUser"}

	response, err := client.Do(query)

	if err != nil {
		return model.Post{}, err
	}

	data := response["updatePost"].(map[string]interface{})
	var updatedPost []model.Post
	mapstructure.Decode(data["post"], &updatedPost)
	return updatedPost[0], nil
}

func (repo PostRepo) LikePost(postId, userId string) error {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			updatePost(input: {filter: {id: "%s"}, set: {likes: {userId: "%s"}}}){
			  post{
				id
			  }
			}
		  }
		  `, postId, userId), Operation: "updatePost"}

	_, err := client.Do(query)

	if err != nil {
		return err
	}

	return nil
}

func (repo PostRepo) UnlikePost(postId, userId string) error {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			updatePost(input: {filter: {id: "%s"}, remove: {likes: {userId: "%s"}}}){
			  post{
				id
			  }
			}
		  }
		  `, postId, userId), Operation: "updatePost"}

	_, err := client.Do(query)

	if err != nil {
		return err
	}

	return nil
}

func (repo PostRepo) DeletePost(postId string) error {

	client := repo.client
	query := dgraph.Request{
		Query: fmt.Sprintf(`mutation {
			updatePost(input: {filter: {id: "%s"}, set: {isDeleted: true}}){
			  post{
				id
			  }
			}
		  }
		  `, postId), Operation: "updatePost"}

	_, err := client.Do(query)

	if err != nil {
		return err
	}

	return nil
}

/**

{ "patch":
  { "filter": {
    "userId": {"eq": "abcd"}
    },
    "set": {
      "books": {"id": "0x15"},
      "posts": {"text": "this is another testing", "book": {"id": "0x15"},
      	"createdAt": "2020-10-23", "updatedAt": "2020-10-23", "isDeleted": false
      }
    }
  }
}


**/
