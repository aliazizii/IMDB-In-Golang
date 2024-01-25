package response

type Err struct {
	Message string `json:"message"`
}

type Movie struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
}

type AllMovies struct {
	Movies []Movie `json:"movies"`
}

type Comment struct {
	ID       int    `json:"id"`
	Author   string `json:"author"`
	Body     string `json:"body"`
	Approved bool   `json:"approved"`
}

type Comments struct {
	Movie    string    `json:"movie"`
	Comments []Comment `json:"comments"`
}

type User struct {
	ID       int
	Username string
	IsAdmin  bool
	JWT      string
}

func CreateErrMessageResponse(s string) Err {
	err := Err{
		Message: s,
	}
	return err
}
