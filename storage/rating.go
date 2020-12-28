package storage

import (
	"context"
	"time"

	"github.com/BarTar213/movies-service/models"
	"github.com/go-pg/pg/v10"
)

func (p *Postgres) AddRating(rating *models.Rating) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		oldRating := &models.Rating{
			UserId:     rating.UserId,
			MovieId:    rating.MovieId,
		}
		err := tx.Model(oldRating).WherePK().Select()
		if err != nil && err != pg.ErrNoRows {
			return err
		}

		if oldRating.Value == rating.Value {
			return nil
		}

		_, err = tx.Model(rating).
			OnConflict("(user_id, movie_id) DO UPDATE").
			Set("value=?value").
			Set("create_date=?create_date").
			Returning(all).
			Insert()
		if err != nil {
			return err
		}

		query := tx.Model((*models.Movie)(nil)).
			Where("id=?", rating.MovieId)

		if oldRating.Value != nil {
			additionalValue := *rating.Value - *oldRating.Value
			query.Set("vote_sum=vote_sum+?", additionalValue)
		} else {
			query.Set("vote_count=vote_count+1")
		}
		_, err = query.Update()
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (p *Postgres) DeleteRating(rating *models.Rating) error {
	_, err := p.db.Model(rating).
		WherePK().
		Delete()

	return err
}
