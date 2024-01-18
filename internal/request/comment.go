package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UpdateComment struct {
	Approved bool `json:"approved"`
}

type Comment struct {
	MovieID     int    `json:"movie_id"`
	CommentBody string `json:"comment_body"`
}

func (uc UpdateComment) Validate() error {
	return validation.ValidateStruct(&uc,
		validation.Field(&uc.Approved, validation.Required),
	)
}

func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.MovieID, validation.Required),
		validation.Field(&c.CommentBody, validation.Required),
	)
}
