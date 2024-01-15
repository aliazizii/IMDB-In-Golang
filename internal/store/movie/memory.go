package movie

import (
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/comment"
)

type MovieInMemory struct {
	Movies   map[int]model.Movie
	movieIDs int
}

var M *MovieInMemory

func NewMovieInMemory() *MovieInMemory {
	M = &MovieInMemory{
		Movies:   make(map[int]model.Movie),
		movieIDs: 0,
	}
	return M
}

func (m *MovieInMemory) isExistMovie(name string, description string) (bool, error) {
	for _, movie := range m.Movies {
		if movie.Name == name && movie.Description == description {
			return true, nil
		}
	}
	return false, nil
}

func (m *MovieInMemory) AddMovie(movie model.Movie) error {
	if ok, _ := m.isExistMovie(movie.Name, movie.Description); ok {
		return DuplicateMovie
	}
	m.movieIDs++
	movie.ID = m.movieIDs
	m.Movies[movie.ID] = movie

	//log

	return nil
}

func (m *MovieInMemory) DeleteMovie(i int) error {
	if _, ok := m.Movies[i]; !ok {
		return MovieNotFound
	}
	delete(m.Movies, i)
	//log
	return nil
}

func (m *MovieInMemory) UpdateMovie(i int, movie model.Movie) error {
	preMovie, ok := m.Movies[i]
	if !ok {
		return MovieNotFound
	}
	movie.ID = i
	movie.Rating = preMovie.Rating
	movie.NVote = preMovie.NVote
	movie.Comments = make([]model.Comment, 0)
	copy(movie.Comments, preMovie.Comments)
	m.Movies[i] = movie
	//log
	return nil
}

func (m *MovieInMemory) AllMovies() ([]model.Movie, error) {
	movies := make([]model.Movie, 0)
	for _, movie := range m.Movies {
		movies = append(movies, movie)
	}
	//log
	return movies, nil
}

func (m *MovieInMemory) Movie(i int) (model.Movie, error) {
	movie, ok := m.Movies[i]
	if !ok {
		return movie, MovieNotFound
	}
	//log
	return movie, nil
}

func (m *MovieInMemory) AllComments(i int) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	_, err := m.Movie(i)
	if err != nil {
		return comments, MovieNotFound
	}
	for _, comment := range comment.M.Comments {
		if comment.MovieID == i {
			comments = append(comments, comment)
		}
	}
	//log
	return comments, nil
}
