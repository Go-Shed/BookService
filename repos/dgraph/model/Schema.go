package model

type User struct {
	UserId    string `json:"userId,omitempty"`
	Username  string `json:"userName,omitempty"`
	UserPhoto string `json:"userPhoto,omitempty"`
	Followers []User `json:"followers,omitempty"`
	Following []User `json:"following,omitempty"`
	Email     string `json:"email,omitempty"`
	Posts     []Post `json:"posts,omitempty"`
	Likes     []Post `json:"likes,omitempty"`
	Books     []Book `json:"books,omitempty"`
}

type Post struct {
	Id             string    `json:"id,omitempty"`
	Text           string    `json:"text,omitempty"`
	CreatedAt      string    `json:"createdAt,omitempty"`
	UpdatedAt      string    `json:"updatedAt,omitempty"`
	LikesAggregate Aggregate `json:"-"`
	Likes          []User    `json:"likes,omitempty"`
	Author         User      `json:"author,omitempty"`
	Book           Book      `json:"book,omitempty"`
	IsDeleted      bool      `json:"isDeleted,omitempty"`
}

type Book struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Posts []Post `json:"posts,omitempty"`
	Users []User `json:"users,omitempty"`
}

type Aggregate struct {
	Count int64 `json:"count,omitempty"`
}

/**
Dgraph schema:

type User {
  userId: String! @id
  userName: String! @search(by: [exact])
  userPhoto: String
  email: String
  followers: [User!]
  following: [User!] @hasInverse(field:followers)
  posts: [Post!] @hasInverse(field:author)
  likes: [Post!] @hasInverse(field:likes)
}



type Post {
  id: ID!
  text: String!
  author: User!
  color: String
  createdAt: DateTime!
  updatedAt: DateTime!
  totalLikes: Int64
  likes: [User!]
}

type Book{
	name: String! @id @search(by: [trigram])
	posts: [Post!]
	users: [User!]
  }
**/
