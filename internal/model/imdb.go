package model

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	Comments    []Comment `json:"comments"`
	NVote       float64   `json:"number_of_votes"`
}

type Vote struct {
	ID      int
	User    User
	Rating  float64
	MovieID int
}

type Comment struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	MovieID   int       `json:"movie_id"`
	Approved  bool      `json:"approved"`
}
type User struct {
	ID       int
	Role     int
	Username string
	Password string // should not save plain text
}
