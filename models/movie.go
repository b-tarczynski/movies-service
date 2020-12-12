package models

import "time"

const emptyStr = ""

type Movie struct {
	Id               int         `json:"id"`
	Adult            bool        `json:"adult"`
	Budget           int64       `json:"budget"`
	BackdropPath     string      `json:"backdrop_path"`
	Homepage         string      `json:"homepage"`
	ImdbId           string      `json:"imdb_id"`
	OriginalLanguage string      `json:"original_language"`
	OriginalTitle    string      `json:"original_title"`
	Overview         string      `json:"overview"`
	Popularity       float32     `json:"popularity"`
	PosterPath       string      `json:"poster_path"`
	ReleaseDate      time.Time        `json:"release_date"`
	Revenue          int64       `json:"revenue"`
	Runtime          int         `json:"runtime"`
	Status           string      `json:"status"`
	Tagline          string      `json:"tagline"`
	Title            string      `json:"title"`
	VoteAverage      float32     `json:"vote_average"`
	VoteCount        int         `json:"vote_count"`
	Countries        []*Country  `json:"production_countries" pg:",many2many:movie_countries"`
	Companies        []*Company  `json:"production_companies" pg:",many2many:movie_companies"`
	Genres           []*Genre    `json:"genres" pg:",many2many:movie_genres"`
	Languages        []*Language `json:"spoken_languages" pg:",many2many:movie_languages"`
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

func (m *Movie) Reset() {
	m.Id = 0
	m.Adult = false
	m.Budget = 0
	m.Homepage = emptyStr
	m.ImdbId = emptyStr
	m.OriginalLanguage = emptyStr
	m.OriginalLanguage = emptyStr
	m.Overview = emptyStr
	m.Popularity = 0
	m.PosterPath = emptyStr
	m.ReleaseDate = time.Time{}
	m.Revenue = 0
	m.Runtime = 0
	m.Status = emptyStr
	m.Tagline = emptyStr
	m.Title = emptyStr
	m.Title = emptyStr
	m.VoteAverage = 0
	m.VoteCount = 0
	m.Companies = nil
	m.Countries = nil
	m.Genres = nil
	m.Languages = nil
}

type Country struct {
	Code string `json:"iso_3166_1" pg:",pk"`
	Name string `json:"name"`
}

type Company struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Language struct {
	IsoCode string `json:"iso_639_1" pg:"iso_639_1,pk"`
	Name    string `json:"name"`
}

type MovieCountry struct {
	MovieId     int
	CountryCode string
}

type MovieCompany struct {
	MovieId   int
	CompanyId int
}

type MovieGenre struct {
	MovieId int
	GenreId int
}

type MovieLanguage struct {
	MovieId int
	IsoCode string `pg:"language_iso_639_1"`
}
