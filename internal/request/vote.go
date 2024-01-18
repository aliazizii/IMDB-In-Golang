package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Vote struct {
	MovieID int `json:"movie_id"`
	Vote    int `json:"vote"`
}

func (v Vote) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.MovieID, validation.Required),
		validation.Field(&v.Vote, validation.Required, validation.In(0, 1, 2, 3, 4, 5)),
	)
}
