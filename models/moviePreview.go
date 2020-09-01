package models

import "time"

type MoviePreview struct {
	tableName   struct{}  `pg:"movies"`
	Id          int       `json:"id"`
	ReleaseDate time.Time `json:"release_date"`
	Title       string    `json:"title"`
	PosterPath  string    `json:"poster_path"`
	VoteAverage float32   `json:"vote_average"`
	Genres      []*Genre  `json:"genres" pg:",many2many:movie_genres"`
}
