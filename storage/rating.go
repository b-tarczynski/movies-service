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
			UserId:  rating.UserId,
			MovieId: rating.MovieId,
		}
		err := tx.Model(oldRating).WherePK().Select()
		if err != nil && err != pg.ErrNoRows {
			return err
		}

		if oldRating.Rating == rating.Rating {
			return nil
		}

		_, err = tx.Model(rating).
			OnConflict("(user_id, movie_id) DO UPDATE").
			Set("rating=?rating").
			Set("create_date=?create_date").
			Returning(all).
			Insert()
		if err != nil {
			return err
		}

		query := tx.Model((*models.Movie)(nil)).
			Where("id=?", rating.MovieId)

		if oldRating.Rating != nil {
			additionalValue := *rating.Rating - *oldRating.Rating
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

func (p *Postgres) ListRatedMovies(userID int) ([]models.MoviePreview, error) {
	movies :=  make([]models.MoviePreview, 0)

	err := p.db.Model((*models.Rating)(nil)).
		Column("m.*").
		Column("rating").
		Where("user_id=?", userID).
		Join("LEFT JOIN movies m ON m.id = rating.movie_id").
		Select(&movies)

	return movies, err
}
