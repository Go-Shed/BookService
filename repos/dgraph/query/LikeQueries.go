package query

import (
	"shed/bookservice/repos/dgraph"
	"shed/bookservice/repos/dgraph/model"
	"time"
)

type LikeRepo struct {
	client dgraph.Dgraph
}

func NewLikeRepo() LikeRepo {
	return LikeRepo{client: dgraph.Dgraph{
		URL:    "https://billowing-night.ap-south-1.aws.cloud.dgraph.io/graphql",
		Secret: "ZTE4YjRhNGEwYTAxNWRiYjMwMGI4YmEyMjc1ODU5ZmI=",
	}}
}

func (repo LikeRepo) UpdateLikeOnPost(postId, userId string) error {

	likeNode := model.Like{
		LikedBy:   model.User{Username: userId},
		LikedOn:   model.Post{Id: postId},
		TimeStamp: time.Now(),
	}

	client := repo.client
	query := dgraph.Request{
		Query: `mutation addLike($patch: [AddLikeInput!]!) {
			addLike(input: $patch) {
				like{
				likedBy{
					username
				}
			}
			}
		}`, Variables: dgraph.Variables{Patch: likeNode}, Operation: "addLike"}

	_, err := client.Do(query)

	if err != nil {
		return err
	}
	return nil
}
