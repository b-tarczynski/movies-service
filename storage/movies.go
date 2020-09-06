package storage

import "github.com/BarTar213/movies-service/models"

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

func (p *Postgres) LikeMovie(userId int, movieId int) error {
	_, err := p.db.ExecOne("INSERT INTO liked_movies (user_id, movie_id) values (?, ?)", userId, movieId)

	return err
}

func (p *Postgres) DeleteMovieLike(userId int, movieId int) error{
	_, err := p.db.ExecOne("DELETE FROM liked_movies WHERE user_id=? AND movie_id=?)", userId, movieId)

	return err
}

func (p *Postgres) AddRecentViewedMovie(userId int, movieId int) error {
	deleteQuery := `
		DELETE
		FROM user_history
		WHERE user_id = ?
		  AND movie_id IN (SELECT movie_id FROM user_history ORDER BY time OFFSET 20)`

	insertQuery := `
		insert into user_history (user_id, movie_id)
		values (?, ?)
		on conflict (user_id, movie_id) do update set time = now()`

	_, err := p.db.Exec(insertQuery, userId, movieId)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(deleteQuery, userId)

	return err
}
