package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
}
