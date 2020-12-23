package models

import "time"

type MoviePreview struct {
	tableName   struct{}  `pg:"movies,discard_unknown_columns"`
	Id          int       `json:"id"`
	PosterPath  string    `json:"poster_path"`
	ReleaseDate time.Time `json:"release_date"`
	VoteAverage float32   `json:"vote_average"`
	Title       string    `json:"title"`
}

type LikedMovie struct {
	MovieId int `pg:",pk"`
	UserId  int `pg:",pk"`
}
