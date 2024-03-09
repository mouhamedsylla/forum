package models

import (
	"forum/orm"
	"time"
)

type User struct {
	orm.Model
	Username string `orm-go:"NOT NULL"`
	Email    string `orm-go:"NOT NULL UNIQUE"`
	Password string `orm-go:"NOT NULL"`
}
type Post struct {
	orm.Model
	Title   string `orm-go:"NOT NULL"`
	Content string `orm-go:"NOT NULL"`
	User_id int    `orm-go:"FOREIGN_KEY:User:Id"`
	Like    int    `json:"nbLike"`
	Dislike int    `json:"nbDislike"`
}
type Comments struct {
	orm.Model
	Comment string `orm-go:"NOT NULL"`
	Post_id int    `orm-go:"FOREIGN_KEY:Post:Id"`
	User_id int    `orm-go:"FOREIGN_KEY:User:Id"`
	Like    int
	Dislike int
}
type Categories struct {
	Name    string `orm-go:"NOT NULL"`
	Id_Post int    `orm-go:"FOREIGN_KEY:Post:Id"`
}

type ReactionPost struct {
	Value   string `json:"value"`
	PostId  int    `orm-go:"FOREIGN_KEY:Post:Id"`
	User_id int    `orm-go:"FOREIGN_KEY:User:Id"`
}

type ReactionComment struct {
	Value   string `json:"value"`
	CommentId  int    `orm-go:"FOREIGN_KEY:Comments:Id"`
	User_id int    `orm-go:"FOREIGN_KEY:User:Id"`
}

type SessionDb struct {
	SessionID string `orm-go:"NOT NULL UNIQUE"`
	User_Id   int    `orm-go:"FOREIGN_KEY:User:Id"`
	TimeOut   string
}

type Message struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Session struct {
	SessionId string    `json:"session_id"`
	User      *User     `json:"user"`
	TimeOut   time.Time `json:"timeout"`
}

type ProfileConnect struct {
	Message Message  `json:"message"`
	Session *Session `json:"session"`
}

type PostID struct {
	Id int
}

type CommentId struct {
	Id int
}

type Action struct {
	Reaction string
	Id   int
}

type PostUI struct {
	User       User
	Post       Post
	Categories []Categories
}
type ErrorData struct {
	ErrorCode int
	Message   string
}
