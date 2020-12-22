package models

type LikedComment struct {
	CommentId int `pg:",pk"`
	UserId    int `pg:",pk"`
}
