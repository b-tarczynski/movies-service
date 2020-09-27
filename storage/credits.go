package storage

import (
	"context"
	"time"

	"github.com/BarTar213/movies-service/models"
	"github.com/go-pg/pg/v10"
)

func (p *Postgres) GetCredits(movieId int, credit *models.Credit) error {
	return p.db.Model(credit).
		Where("movie_id = ?", movieId).
		Relation("Cast").
		Relation("Crew").
		Select()
}

func (p *Postgres) AddCredits(credit *models.Credit) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		_, err := tx.Model(credit).Insert()
		if err != nil {
			return err
		}
		_, err = tx.Model(&credit.Cast).Insert()
		if err != nil {
			return err
		}
		_, err = tx.Model(&credit.Crew).Insert()
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
