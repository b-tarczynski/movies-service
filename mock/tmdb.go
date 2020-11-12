package mock

import (
	"net/http"

	"github.com/BarTar213/movies-service/models"
)

type Tmdb struct {
	GetCreditsErr    bool
	GetCreditsStatus int
}

func (t *Tmdb) GetCredits(movieId int, credit *models.Credit) (int, error) {
	if t.GetCreditsErr {
		return t.GetCreditsStatus, exampleErr
	}
	return http.StatusOK, nil
}
