package mock

import (
	"errors"

	"github.com/BarTar213/movies-service/models"
	"github.com/go-pg/pg/v10"
)

var exampleErr = errors.New("example error")

type Storage struct {
	GetMovieErr             bool
	ListMoviesErr           bool
	LikeMovieErr            bool
	DeleteMovieLikeErr      bool
	AddRecentViewedMovieErr bool
	ListLikedMoviesErr      bool
	CheckLikedErr           bool

	GetMovieCommentsErr  bool
	AddMovieCommentErr   bool
	UpdateCommentErr     bool
	LikeCommentErr       bool
	DeleteCommentLikeErr bool
	DeleteCommentErr     bool

	GetCreditsErr         bool
	GetCreditsNotFoundErr bool
	AddCreditsErr         bool
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

func (s *Storage) ListMovies(title string, params *models.PaginationParams) ([]models.Movie, error) {
	if s.ListMoviesErr {
		return nil, exampleErr
	}
	return []models.Movie{}, nil
}

func (s *Storage) LikeMovie(userId int, movieId int) error {
	if s.LikeMovieErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) DeleteMovieLike(userId int, movieId int) error {
	if s.DeleteMovieLikeErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) AddRecentViewedMovie(userId int, movieId int) error {
	if s.AddRecentViewedMovieErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) LikeComment(userId int, commentId int) error {
	if s.LikeCommentErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) DeleteCommentLike(userId int, commentId int) error {
	if s.DeleteCommentLikeErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) GetCredits(movieId int, credit *models.Credit) error {
	if s.GetCreditsErr {
		return exampleErr
	}
	if s.GetCreditsNotFoundErr {
		return pg.ErrNoRows
	}
	return nil
}

func (s *Storage) AddCredits(credit *models.Credit) error {
	if s.AddCreditsErr {
		return exampleErr
	}
	return nil
}

func (s *Storage) ListLikedMovies(userId int, params *models.PaginationParams) ([]models.MoviePreview, error) {
	if s.ListLikedMoviesErr {
		return nil, exampleErr
	}
	return []models.MoviePreview{}, nil
}

func (s *Storage) CheckLiked(likedMovie *models.LikedMovie) (bool, error) {
	if s.CheckLikedErr {
		return false, exampleErr
	}
	return true, nil
}
