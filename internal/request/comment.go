package request

type UpdateComment struct {
	Approved bool `json:"approved"`
}

type Comment struct {
	MovieID     int    `json:"movie_id"`
	CommentBody string `json:"comment_body"`
}
