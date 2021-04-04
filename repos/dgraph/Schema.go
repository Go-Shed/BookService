package dgraph

import "time"

type DgraphResponse struct {
	Data interface{} `json:"data"`
}

type User struct {
	Username  string `json:"username"`
	Followers []User `json:"followers,omitempty"`
	Following []User `json:"following,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	Posts     []Post `json:"posts,omitempty"`
	Likes     []Like `json:"likes,omitempty"`
}

type Post struct {
	Id            string    `json:"id"`
	Title         string    `json:"title,omitempty"`
	Text          string    `json:"text,omitempty"`
	Tags          string    `json:"tags,omitempty"`
	DatePublished time.Time `json:"datePublished,omitempty"`
	TotalLikes    int64     `json:"totalLikes,omitempty"`
	Likes         []Like    `json:"likes,omitempty"`
}

type Like struct {
	Id        string    `json:"id"`
	TimeStamp time.Time `json:"timestamp,omitempty"`
	LikedOn   Post      `json:"likedOn,omitempty"`
	LikedBy   User      `json:"likedBy,omitempty"`
}
