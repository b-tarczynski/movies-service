package storage

import (
	"context"
	"time"

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
	DeleteMovieLike(userId int, movieId int) error
	AddRecentViewedMovie(userId int, movieId int) error

	GetMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error)
	LikeComment(userId int, commentId int) error
	DeleteCommentLike(userId int, commentId int) error
	AddMovieComment(comment *models.Comment) error
	UpdateComment(comment *models.Comment) error
	DeleteComment(comment *models.Comment) error

	GetCredits(movieId int, credit *models.Credit) error
	AddCredits(credit *models.Credit) error
}

func NewPostgres(config *config.Postgres) (Storage, error) {
	db := pg.Connect(&pg.Options{
		Addr:        config.Address,
		User:        config.User,
		Password:    config.Password,
		Database:    config.Database,
		DialTimeout: 5 * time.Second,
	})

	err := db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}
