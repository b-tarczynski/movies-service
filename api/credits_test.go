package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/mock"
	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/BarTar213/movies-service/tmdb"
	"github.com/gin-gonic/gin"
)

func TestMovieHandlers_GetCredits(t *testing.T) {
	type fields struct {
		storage    storage.Storage
		conf       *config.Config
		tmdbClient tmdb.Client
		logger     *log.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		movieId    string
		wantStatus int
	}{
		{
			name: "positive_get_credits",
			fields: fields{
				storage:    &mock.Storage{},
				conf:       &config.Config{},
				tmdbClient: &mock.Tmdb{},
				logger:     logger,
			},
			movieId:    validId,
			wantStatus: http.StatusOK,
		},
		{
			name: "negative_get_credits_invalid_movie_id",
			fields: fields{
				storage:    &mock.Storage{},
				conf:       &config.Config{},
				tmdbClient: &mock.Tmdb{},
				logger:     logger,
			},
			movieId:    invalidId,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_get_credits_storage_error",
			fields: fields{
				storage: &mock.Storage{
					GetCreditsErr: true,
				},
				conf:       &config.Config{},
				tmdbClient: &mock.Tmdb{},
				logger:     logger,
			},
			movieId:    validId,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "positive_get_credits_storage_error",
			fields: fields{
				storage: &mock.Storage{
					GetCreditsNotFoundErr: true,
				},
				conf:       &config.Config{},
				tmdbClient: &mock.Tmdb{},
				logger:     logger,
			},
			movieId:    validId,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
				WithTmdbClient(tt.fields.tmdbClient),
			)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/movies/%s/credits", tt.movieId)
			req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestMovieHandlers_AddCredits(t *testing.T) {
	type fields struct {
		storage storage.Storage
		tmdb    tmdb.Client
		logger  *log.Logger
	}
	type args struct {
		credits *models.Credit
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive_add_credits",
			fields: fields{
				storage: &mock.Storage{},
				tmdb:    &mock.Tmdb{},
				logger:  logger,
			},
			args: args{
				credits: &models.Credit{
					Id:      1,
					MovieId: 1,
					Cast: []*models.Cast{
						{},
					},
					Crew: []*models.Crew{
						{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "negative_add_credits_storage_error",
			fields: fields{
				storage: &mock.Storage{
					AddCreditsErr: true,
				},
				tmdb:   &mock.Tmdb{},
				logger: logger,
			},
			args: args{
				credits: &models.Credit{
					Id:      1,
					MovieId: 1,
					Cast: []*models.Cast{
						{},
					},
					Crew: []*models.Crew{
						{},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MovieHandlers{
				storage: tt.fields.storage,
				tmdb:    tt.fields.tmdb,
				logger:  tt.fields.logger,
			}
			if err := h.AddCredits(tt.args.credits); (err != nil) != tt.wantErr {
				t.Errorf("AddCredits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
