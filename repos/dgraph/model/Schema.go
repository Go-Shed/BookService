package model

import "time"

type User struct {
	Username  string `json:"username,omitempty"`
	Followers []User `json:"followers,omitempty"`
	Following []User `json:"following,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	Posts     []Post `json:"posts,omitempty"`
	Likes     []Like `json:"likes,omitempty"`
}

type Post struct {
	Id            string    `json:"id,omitempty"`
	Title         string    `json:"title,omitempty"`
	Text          string    `json:"text,omitempty"`
	Tags          string    `json:"tags,omitempty"`
	DatePublished time.Time `json:"datePublished,omitempty"`
	TotalLikes    int64     `json:"totalLikes,omitempty"`
	Likes         []Like    `json:"likes,omitempty"`
	Author        User      `json:"author,omitempty"`
}

type Like struct {
	Id        string    `json:"id,omitempty"`
	TimeStamp time.Time `json:"timestamp,omitempty"`
	LikedOn   Post      `json:"likedOn,omitempty"`
	LikedBy   User      `json:"likedBy,omitempty"`
}
