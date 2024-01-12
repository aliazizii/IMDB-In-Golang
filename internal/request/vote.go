package request

type Vote struct {
	MovieID int     `json:"movie_id"`
	Vote    float64 `json:"vote"`
}
