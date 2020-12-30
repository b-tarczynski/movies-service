package models

import "time"

type TmdbMovie struct {
	tableName        struct{}    `pg:"movies,discard_unknown_columns,alias:movie"`
	Id               int         `json:"id"`
	Adult            bool        `json:"adult" pg:",use_zero"`
	Budget           int64       `json:"budget" pg:",use_zero"`
	BackdropPath     string      `json:"backdrop_path" pg:",use_zero"`
	Homepage         string      `json:"homepage" pg:",use_zero"`
	ImdbId           string      `json:"imdb_id" pg:",use_zero"`
	OriginalLanguage string      `json:"original_language" pg:",use_zero"`
	OriginalTitle    string      `json:"original_title" pg:",use_zero"`
	Overview         string      `json:"overview" pg:",use_zero"`
	Popularity       float32     `json:"popularity" pg:",use_zero"`
	PosterPath       string      `json:"poster_path" pg:",use_zero"`
	ReleaseDate      Time        `json:"release_date"`
	Revenue          int64       `json:"revenue" pg:",use_zero"`
	Runtime          int         `json:"runtime" pg:",use_zero"`
	Status           string      `json:"status" pg:",use_zero"`
	Tagline          string      `json:"tagline" pg:",use_zero"`
	Title            string      `json:"title" pg:",use_zero"`
	VoteAverage      float32     `json:"vote_average" pg:",use_zero"`
	VoteCount        int         `json:"vote_count" pg:",use_zero"`
	Countries        []*Country  `json:"production_countries" pg:",many2many:movie_countries,fk:movie_id"`
	Companies        []*Company  `json:"production_companies" pg:",many2many:movie_companies,fk:movie_id"`
	Genres           []*Genre    `json:"genres" pg:",many2many:movie_genres,fk:movie_id"`
	Languages        []*Language `json:"spoken_languages" pg:",many2many:movie_languages,fk:movie_id"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	if len(string(data)) <= 2 {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	tt, err := time.Parse(`"`+"2006-01-02"+`"`, string(data))
	*t = Time{tt}
	return err
}
