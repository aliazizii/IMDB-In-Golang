package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Movie struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m Movie) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Description, validation.Required),
	)
}
