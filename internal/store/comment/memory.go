package comment

import (
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/movie"
)

type CommentInMemory struct {
	Comments   map[int]model.Comment
	commentIDs int
}

var M *CommentInMemory

func NewCommentInMemory() *CommentInMemory {
	M = &CommentInMemory{
		Comments:   make(map[int]model.Comment),
		commentIDs: 0,
	}
	return M
}

func (m *CommentInMemory) UpdateComment(i int, approved bool) error {
	comment, ok := m.Comments[i]
	if !ok {
		return CommentNotFound
	}
	comment.Approved = approved
	m.Comments[i] = comment
	for j := range movie.M.Movies[comment.MovieID].Comments {
		if movie.M.Movies[comment.MovieID].Comments[j].ID == i {
			movie.M.Movies[comment.MovieID].Comments[j].Approved = approved
		}
	}
	return nil
}

func (m *CommentInMemory) DeleteComment(i int) error {
	comment, ok := m.Comments[i]
	if !ok {
		return CommentNotFound
	}
	delete(m.Comments, i)
	j := 0
	for j = range movie.M.Movies[comment.MovieID].Comments {
		if movie.M.Movies[comment.MovieID].Comments[j].ID == i {
			break
		}
	}
	// Should be tested
	movie.M.Movies[comment.MovieID].Comments[j] = movie.M.Movies[comment.MovieID].Comments[len(movie.M.Movies[comment.MovieID].Comments)-1]
	copy(movie.M.Movies[comment.MovieID].Comments, movie.M.Movies[comment.MovieID].Comments[:len(movie.M.Movies[comment.MovieID].Comments)-1])
	return nil
}

func (m *CommentInMemory) Comment(comment model.Comment) error {
	m.commentIDs++
	comment.ID = m.commentIDs
	m.Comments[comment.ID] = comment
	tempMovie := movie.M.Movies[comment.MovieID]
	tempMovie.Comments = append(tempMovie.Comments, comment)
	movie.M.Movies[comment.MovieID] = tempMovie
	//log
	return nil
}
