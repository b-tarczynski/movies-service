package storage

import (
	"github.com/BarTar213/movies-service/models"
)

func (p *Postgres) ListMovieComments(movieId int, params *models.PaginationParams) ([]models.Comment, error) {
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

func (p *Postgres) ListLikedCommentsForMovie(movieID, userID int) ([]int, error) {
	ids := make([]int, 0)

	err := p.db.Model((*models.LikedComment)(nil)).
		Column("comment_id").
		Where("c.movie_id = ?", movieID).
		Where("liked_comment.user_id = ?", userID).
		Join("LEFT JOIN comments c ON liked_comment.comment_id = c.id").
		Select(&ids)

	return ids, err
}

func (p *Postgres) LikeComment(userId int, commentId int, comment *models.Comment) error {
	_, err := p.db.ExecOne("INSERT INTO liked_comments (user_id, comment_id) VALUES (?, ?)", userId, commentId)
	if err != nil {
		return err
	}

	comment.Id = commentId
	err = p.db.Model(comment).
		WherePK().
		Select()

	return err
}

func (p *Postgres) DeleteCommentLike(userId int, commentId int) error {
	_, err := p.db.ExecOne("DELETE FROM liked_comments WHERE user_id=? AND comment_id=?", userId, commentId)

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
