package store

import "github.com/aliazizii/IMDB-In-Golang/internal/model"

type IMDBInMemory struct {
	movies   map[int]model.Movie
	votes    map[int]model.Vote
	comments map[int]model.Comment
	users    map[int]model.User
}

func NewIMDBInMemory() *IMDBInMemory {
	return &IMDBInMemory{
		movies:   make(map[int]model.Movie),
		votes:    make(map[int]model.Vote),
		comments: make(map[int]model.Comment),
		users:    make(map[int]model.User),
	}
}

func (m *IMDBInMemory) AddMovie(movie model.Movie) error {
	if _, ok := m.movies[movie.ID]; ok {
		return DuplicateMovie
	}
	m.movies[movie.ID] = movie

	//log

	return nil
}

func (m *IMDBInMemory) DeleteMovie(i int) error {
	if _, ok := m.movies[i]; !ok {
		return MovieNotFound
	}
	delete(m.movies, i)
	//log
	return nil
}

func (m *IMDBInMemory) UpdateMovie(i int, movie model.Movie) error {
	if _, ok := m.movies[i]; !ok {
		return MovieNotFound
	}
	m.movies[i] = movie
	//log
	return nil
}

func (m *IMDBInMemory) Vote(vote model.Vote) error {
	if _, ok := m.votes[vote.ID]; ok {
		return DuplicateVote
	}
	m.votes[vote.ID] = vote
	tempMovie := m.movies[vote.MovieID]
	tempMovie.Rating = (tempMovie.Rating*tempMovie.NVote + vote.Rating) / (tempMovie.NVote + 1)
	tempMovie.NVote++
	m.movies[vote.MovieID] = tempMovie
	//log
	return nil
}

func (m *IMDBInMemory) Comment(comment model.Comment) error {
	m.comments[comment.ID] = comment
	tempMovie := m.movies[comment.MovieID]
	tempMovie.Comments = append(tempMovie.Comments, comment)
	m.movies[comment.MovieID] = tempMovie
	//log
	return nil
}

func (m *IMDBInMemory) AllMovies() ([]model.Movie, error) {
	movies := make([]model.Movie, 0)
	for _, movie := range m.movies {
		movies = append(movies, movie)
	}
	//log
	return movies, nil
}

func (m *IMDBInMemory) Movie(i int) (model.Movie, error) {
	movie, ok := m.movies[i]
	if !ok {
		return movie, MovieNotFound
	}
	//log
	return movie, nil
}

func (m *IMDBInMemory) AllComments() ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	for _, comment := range m.comments {
		comments = append(comments, comment)
	}
	//log
	return comments, nil
}
