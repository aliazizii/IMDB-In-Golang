package store

import (
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
)

type IMDBInMemory struct {
	movies     map[int]model.Movie
	votes      map[int]model.Vote
	comments   map[int]model.Comment
	users      map[int]model.User
	movieIDs   int
	voteIDs    int
	commentIDs int
	userIDs    int
}

func NewIMDBInMemory() *IMDBInMemory {
	return &IMDBInMemory{
		movies:     make(map[int]model.Movie),
		votes:      make(map[int]model.Vote),
		comments:   make(map[int]model.Comment),
		users:      make(map[int]model.User),
		movieIDs:   0,
		voteIDs:    0,
		commentIDs: 0,
		userIDs:    0,
	}
}

// Admin

func (m *IMDBInMemory) AddMovie(movie model.Movie) error {
	if ok, _ := m.isExistMovie(movie.Name, movie.Description); ok {
		return DuplicateMovie
	}
	m.movieIDs++
	movie.ID = m.movieIDs
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
	preMovie, ok := m.movies[i]
	if !ok {
		return MovieNotFound
	}
	movie.ID = i
	movie.Rating = preMovie.Rating
	movie.NVote = preMovie.NVote
	movie.Comments = make([]model.Comment, 0)
	copy(movie.Comments, preMovie.Comments)
	m.movies[i] = movie
	//log
	return nil
}
func (m *IMDBInMemory) UpdateComment(i int, approved bool) error {
	comment, ok := m.comments[i]
	if !ok {
		return CommentNotFound
	}
	comment.Approved = approved
	m.comments[i] = comment
	for j := range m.movies[comment.MovieID].Comments {
		if m.movies[comment.MovieID].Comments[j].ID == i {
			m.movies[comment.MovieID].Comments[j].Approved = approved
		}
	}
	return nil
}

func (m *IMDBInMemory) DeleteComment(i int) error {
	comment, ok := m.comments[i]
	if !ok {
		return CommentNotFound
	}
	delete(m.comments, i)
	j := 0
	for j = range m.movies[comment.MovieID].Comments {
		if m.movies[comment.MovieID].Comments[j].ID == i {
			break
		}
	}
	// Should be tested
	m.movies[comment.MovieID].Comments[j] = m.movies[comment.MovieID].Comments[len(m.movies[comment.MovieID].Comments)-1]
	copy(m.movies[comment.MovieID].Comments, m.movies[comment.MovieID].Comments[:len(m.movies[comment.MovieID].Comments)-1])
	return nil
}

// User

func (m *IMDBInMemory) Vote(vote model.Vote) error {
	if ok, _ := m.isExistVote(vote.MovieID, vote.User.ID); ok {
		return DuplicateVote
	}
	m.voteIDs++
	vote.ID = m.voteIDs
	m.votes[vote.ID] = vote

	tempMovie := m.movies[vote.MovieID]
	tempMovie.Rating = (tempMovie.Rating*tempMovie.NVote + vote.Rating) / (tempMovie.NVote + 1)
	tempMovie.NVote++
	m.movies[vote.MovieID] = tempMovie
	//log
	return nil
}

func (m *IMDBInMemory) Comment(comment model.Comment) error {
	m.commentIDs++
	comment.ID = m.commentIDs
	m.comments[comment.ID] = comment
	tempMovie := m.movies[comment.MovieID]
	tempMovie.Comments = append(tempMovie.Comments, comment)
	m.movies[comment.MovieID] = tempMovie
	//log
	return nil
}

// Public

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

func (m *IMDBInMemory) AllComments(i int) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	_, err := m.Movie(i)
	if err != nil {
		return comments, MovieNotFound
	}
	for _, comment := range m.comments {
		if comment.MovieID == i {
			comments = append(comments, comment)
		}
	}
	//log
	return comments, nil
}

func (m *IMDBInMemory) isExistMovie(name string, description string) (bool, error) {
	for _, movie := range m.movies {
		if movie.Name == name && movie.Description == description {
			return true, nil
		}
	}
	return false, nil
}

func (m *IMDBInMemory) isExistVote(movieID int, userID int) (bool, error) {
	for _, vote := range m.votes {
		if movieID == vote.MovieID && userID == vote.User.ID {
			return true, nil
		}
	}
	return false, nil
}
