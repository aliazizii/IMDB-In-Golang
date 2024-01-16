package response

type Err struct {
	message string `json:"message"`
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
	ID     int    `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
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

func CreateErrMessageResponse(s string) *Err {
	return &Err{
		message: s,
	}
}
