package model

import "time"

type Movie struct {
	ID          int
	Name        string
	Description string
	Rating      float64
	Comments    []Comment
}

type Vote struct {
	ID      int
	User    User
	Rating  int
	MovieID int
}

type Comment struct {
	ID        int
	User      User
	Text      string
	CreatedAt time.Time
	MovieID   int
	//Approved bool
}
type User struct {
	ID       int
	Role     int
	Username string
	Password string
}
