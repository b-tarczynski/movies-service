package storage

import (
	"context"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/models"
	"github.com/go-pg/pg/v10"
)

type Postgres struct {
	db *pg.DB
}

type Storage interface {
	GetMovie(movie *models.Movie) error
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
		Relation("Genres").
		Relation("Countries").
		Relation("Companies").
		Relation("Languages").
		Select()

	return err
}
