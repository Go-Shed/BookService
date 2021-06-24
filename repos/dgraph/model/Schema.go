package model

type User struct {
	UserId             string    `json:"userId,omitempty"`
	Username           string    `json:"userName,omitempty"`
	UserPhoto          string    `json:"userPhoto,omitempty"`
	Followers          []User    `json:"followers,omitempty"`
	Following          []User    `json:"following,omitempty"`
	Email              string    `json:"email,omitempty"`
	Posts              []Post    `json:"posts,omitempty"`
	Likes              []Post    `json:"likes,omitempty"`
	Books              []Book    `json:"books,omitempty"`
	FollowersAggregate Aggregate `json:"-"`
	FollowingAggregate Aggregate `json:"-"`
	Comments           []Comment `json:"comments,omitempty"`
	FCMToken           string    `json:"fcmToken,omitempty"`
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
	Comments       []Comment `json:"comments,omitempty"`
}

type Book struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Posts []Post `json:"posts,omitempty"`
	Users []User `json:"users,omitempty"`
}

type Comment struct {
	Id        string `json:"id,omitempty"`
	Text      string `json:"text,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Post      Post   `json:"post,omitempty"`
	User      User   `json:"user,omitempty"`
}

type Aggregate struct {
	Count int64 `json:"count,omitempty"`
}

/**
Dgraph schema:

V1
type User {
  userId: String! @id
  userName: String! @search(by: [exact, trigram])
  userPhoto: String
  email: String
  followers: [User!]
  following: [User!] @hasInverse(field:followers)
  posts: [Post!] @hasInverse(field:author)
  likes: [Post!] @hasInverse(field:likes)
  books:[Book!] @hasInverse(field:users)
}



type Post {
  id: ID!
  text: String!
  author: User!
  createdAt: DateTime!
  updatedAt: DateTime!
  likes: [User!]
  book: Book! @hasInverse(field:posts)
  isDeleted: Boolean!
}

type Book{
 	id: ID!
	name: String! @search(by: [exact])
	posts: [Post!]
	users: [User!]
  }



 V2


 type User {
  userId: String! @id
  userName: String! @search(by: [exact, trigram])
  userPhoto: String
  email: String
  followers: [User!]
  following: [User!] @hasInverse(field:followers)
  posts: [Post!] @hasInverse(field:author)
  likes: [Post!] @hasInverse(field:likes)
  books:[Book!] @hasInverse(field:users)
  comments: [Comment!] @hasInverse(field:user)
}



type Post {
  id: ID!
  text: String!
  author: User!
  createdAt: DateTime!
  updatedAt: DateTime!
  likes: [User!]
  book: Book! @hasInverse(field:posts)
  isDeleted: Boolean!
  comments: [Comment!] @hasInverse(field:post)
}

type Book{
  id: ID!
	name: String! @search(by: [exact])
	posts: [Post!]
	users: [User!]
  }

type Comment {
  id: ID!
  text: String!
  post: Post!
  user: User!
	createdAt: DateTime!
  updatedAt: DateTime!
}
**/
