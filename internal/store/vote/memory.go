package vote

import (
	"github.com/aliazizii/IMDB-In-Golang/internal/model"
	"github.com/aliazizii/IMDB-In-Golang/internal/store/movie"
)

type VoteInMemory struct {
	votes   map[int]model.Vote
	voteIDs int
}

var M *VoteInMemory

func NewMovieInMemory() *VoteInMemory {
	M = &VoteInMemory{
		votes:   make(map[int]model.Vote),
		voteIDs: 0,
	}
	return M
}

func (m *VoteInMemory) isExistVote(movieID int, userUsername string) (bool, error) {
	for _, vote := range m.votes {
		if movieID == vote.MovieID && userUsername == vote.User.Username {
			return true, nil
		}
	}
	return false, nil
}

func (m *VoteInMemory) Vote(vote model.Vote) error {
	if ok, _ := m.isExistVote(vote.MovieID, vote.User.Username); ok {
		return DuplicateVote
	}
	m.voteIDs++
	vote.ID = m.voteIDs
	m.votes[vote.ID] = vote

	tempMovie := movie.M.Movies[vote.MovieID]
	tempMovie.Rating = (tempMovie.Rating*tempMovie.NVote + vote.Rating) / (tempMovie.NVote + 1)
	tempMovie.NVote++
	movie.M.Movies[vote.MovieID] = tempMovie
	//log
	return nil
}
