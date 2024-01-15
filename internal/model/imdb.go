package model

import "time"

type Movie struct {
	ID          int       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	Comments    []Comment `json:"comments"`
	NVote       float64   `json:"number_of_votes"`
}

type Vote struct {
	ID           int    `gorm:"primaryKey;AUTO_INCREMENT"`
	UserUsername string `json:"user_username"`
	User         User   `gorm:"foreignKey:UserUsername"`
	Rating       float64
	MovieID      int
}

type Comment struct {
	ID             int       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	User           User      `json:"user" gorm:"foreignKey:AuthorUsername"`
	AuthorUsername string    `json:"author_username"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"createdAt"`
	MovieID        int       `json:"movie_id"`
	Approved       bool      `json:"approved"`
}
type User struct {
	ID       int `gorm:"primaryKey;AUTO_INCREMENT"`
	Role     int
	Username string `gorm:"unique"`
	Password string // should not save plain text
}
