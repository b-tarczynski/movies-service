package storage

import (
	"context"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/models"
	"github.com/go-pg/pg/v10"
)

const (
	all = "*"
	genres = "Genres"
)

type Postgres struct {
	db *pg.DB
}

type Storage interface {
	GetMovie(movie *models.Movie) error
	ListMovies(title string, params *models.PaginationParams) ([]models.Movie, error)
	LikeMovie(userId int, movieId int) error

	GetMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error)
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

func (p *Postgres) GetMovie(movie *models.Movie) error {
	err := p.db.Model(movie).
		WherePK().
		Relation(genres).
		Relation("Countries").
		Relation("Companies").
		Relation("Languages").
		Select()

	return err
}

func (p *Postgres) ListMovies(title string, params *models.PaginationParams) ([]models.Movie, error) {
	movies := make([]models.Movie, 0)
	err := p.db.Model(&movies).
		Where("title like ?", title).
		Relation(genres).
		Order(params.OrderBy).
		Offset(params.Offset).
		Limit(params.Limit).
		Select()

	return movies, err
}

func (p *Postgres) GetMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error) {
	comments := make([]models.Comment, 0)
	err := p.db.Model(&comments).
		Where("movie_id = ?", movieId).
		Order(params.OrderBy).
		Offset(params.Offset).
		Limit(params.Limit).
		Select()

	return comments, err
}

func (p *Postgres) LikeMovie(userId int, movieId int) error {
	_, err := p.db.ExecOne("INSERT INTO liked_movies (user_id, movie_id) values (?, ?)", userId, movieId)

	return err
}

func (p *Postgres) AddMovieComment(comment *models.Comment) error {
	_, err := p.db.Model(comment).Returning(all).Insert()

	return err
}

func (p *Postgres) UpdateComment(comment *models.Comment) error {
	_, err := p.db.Model(comment).
		WherePK().
		Where("user_id = ?user_id").
		Set("content = ?content, update_date = ?update_date").
		Returning(all).
		Update()

	return err
}

func (p *Postgres) DeleteComment(comment *models.Comment) error {
	_, err := p.db.Model(comment).
		WherePK().
		Where("user_id = ?user_id").
		Delete()

	return err
}
