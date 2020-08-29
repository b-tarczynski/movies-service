package models

import "time"

type Movie struct {
	Id               bool       `json:"id"`
	Adult            bool       `json:"adult"`
	Budget           int        `json:"budget"`
	Homepage         string     `json:"homepage"`
	ImdbId           string     `json:"imdb_id"`
	OriginalLanguage string     `json:"original_language"`
	OriginalTitle    string     `json:"original_title"`
	Overview         string     `json:"overview"`
	Popularity       float32    `json:"popularity"`
	PosterPath       string     `json:"poster_path"`
	ReleaseDate      time.Time  `json:"release_date"`
	Revenue          int        `json:"revenue"`
	Runtime          int        `json:"runtime"`
	Status           string     `json:"status"`
	Tagline          string     `json:"tagline"`
	Title            string     `json:"title"`
	VoteAverage      float32    `json:"vote_average"`
	VoteCount        int        `json:"vote_count"`
	Countries        []Country  `json:"countries" pg:"many2many:movie_countries"`
	Companies        []Company  `json:"companies" pg:"many2many:movie_companies"`
	Genres           []Genre    `json:"genres" pg:"many2many:movie_genres"`
	Languages        []Language `json:"languages" pg:"many2many:movie_languages"`
}

type Country struct {
	Code string
	Name string
}

type Company struct {
	Id   int
	Name string
}

type Genre struct {
	Id   int
	Name string
}

type Language struct {
	IsoCode string
	Name    string
}
