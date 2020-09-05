package storage

import (
	"context"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/models"
	"github.com/go-pg/pg/v10"
)

const (
	all    = "*"
	genres = "Genres"
)

type Postgres struct {
	db *pg.DB
}

type Storage interface {
	GetMovie(movie *models.Movie) error
	ListMovies(title string, params *models.PaginationParams) ([]models.Movie, error)
	LikeMovie(userId int, movieId int) error
	RemoveMovieLike(userId int, movieId int) error
	AddRecentViewedMovie(userId int, movieId int) error

	GetMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error)
	LikeComment(userId int, commentId int) error
	RemoveCommentLike(userId int, commentId int) error
	AddMovieComment(comment *models.Comment) error
	UpdateComment(comment *models.Comment) error
	DeleteComment(comment *models.Comment) error
}

func NewPostgres(config *config.Postgres) (Storage, error) {
	db := pg.Connect(&pg.Options{
		Addr:     config.Address,
		User:     config.User,
		Password: config.Password,
		Database: config.Database,
	})

	err := db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}
