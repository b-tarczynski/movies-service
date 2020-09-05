package models

import "time"

type Comment struct {
	tableName  struct{}  `pg:",discard_unknown_columns"`
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	MovieId    int       `json:"movie_id"`
	UpdateDate time.Time `json:"update_date"`
	CreateDate time.Time `json:"create_date"`
	Content    string    `json:"content" binding:"required"`
	Likes      int       `json:"likes"`
}
