package storage

import "github.com/BarTar213/movies-service/models"

func (p *Postgres) GetMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error) {
	comments := make([]models.Comment, 0)

	err := p.db.Model(&comments).
		Where("movie_id = ?", movieId).
		ColumnExpr("comment.*").
		ColumnExpr("count(lc.*) AS likes").
		Join("LEFT JOIN liked_comments lc ON comment.id = lc.comment_id").
		Group("comment.id").
		Order(params.OrderBy).
		Offset(params.Offset).
		Limit(params.Limit).
		Select()

	return comments, err
}

func (p *Postgres) LikeComment(userId int, commentId int) error {
	_, err := p.db.ExecOne("INSERT INTO liked_comments (user_id, comment_id) values (?, ?)", userId, commentId)

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
