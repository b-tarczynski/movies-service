package mock

import (
	"errors"

	"github.com/BarTar213/movies-service/models"
)

var exampleErr = errors.New("example error")

type Storage struct {
	GetMovieErr         bool
	GetMovieCommentsErr bool
	AddMovieCommentErr  bool
	UpdateCommentErr    bool
	DeleteCommentErr    bool
}

func (s *Storage) GetMovie(movie *models.Movie) error {
	if s.GetMovieErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) GetMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error) {
	if s.GetMovieCommentsErr {
		return nil, exampleErr
	}
	return []models.Comment{}, nil
}

func (s *Storage) AddMovieComment(comment *models.Comment) error {
	if s.AddMovieCommentErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) UpdateComment(comment *models.Comment) error {
	if s.UpdateCommentErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) DeleteComment(comment *models.Comment) error {
	if s.DeleteCommentErr {
		return exampleErr
	}
	return nil
}