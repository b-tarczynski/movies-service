package mock

import (
	"net/http"

	"github.com/BarTar213/movies-service/models"
)

type Tmdb struct {
	GetCreditsErr    bool
	GetCreditsStatus int

	GetTrendingMoviesErr    bool
	GetTrendingMoviesStatus int
}

func (t *Tmdb) GetCredits(movieId int, credit *models.Credit) (int, error) {
	if t.GetCreditsErr {
		return t.GetCreditsStatus, exampleErr
	}
	return http.StatusOK, nil
}

func (t *Tmdb) GetTrendingMovies() ([]models.TmdbMovie, int, error) {
	if t.GetTrendingMoviesErr {
		return nil,t.GetTrendingMoviesStatus, exampleErr
	}
	return []models.TmdbMovie{}, http.StatusOK, nil
}

