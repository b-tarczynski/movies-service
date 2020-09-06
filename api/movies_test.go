package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/mock"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
)

const (
	accountHeaderId = "X-Account-Id"
	accountHeader   = "X-Account"

	validId     = "2"
	invalidId   = "2ab"
	accountName = "account"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func TestNewMovieHandlers(t *testing.T) {
	type args struct {
		postgres storage.Storage
		logger   *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *MovieHandlers
	}{
		{
			name: "positiveNewMovieHandlers",
			args: args{
				postgres: &mock.Storage{},
				logger:   &log.Logger{},
			},
			want: &MovieHandlers{
				storage: &mock.Storage{},
				logger:  &log.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMovieHandlers(tt.args.postgres, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMovieHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieHandlers_AddRecentViewedMovie(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		accountId   int
		accountName string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
			)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/movies", nil)
			setAccountHeaders(req, tt.accountId, tt.accountName)

			a.Router.ServeHTTP(w, req)
		})
	}
}

func TestMovieHandlers_GetMovie(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		accountId   int
		accountName string
		movieId     string
		wantStatus  int
	}{
		{
			name: "positive_get_movie",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			wantStatus:  http.StatusOK,
		},
		{
			name: "positive_get_movie_without_account_info",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			movieId:    validId,
			wantStatus: http.StatusOK,
		},
		{
			name: "positive_get_movie_add_recent_viewed_error_pass",
			fields: fields{
				storage: &mock.Storage{
					AddRecentViewedMovieErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			wantStatus:  http.StatusOK,
		},
		{
			name: "negative_get_movie_invalid_movie_id",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     invalidId,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_get_movie_storage_error",
			fields: fields{
				storage: &mock.Storage{
					GetMovieErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			wantStatus:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
			)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/movies/%s", tt.movieId)
			req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)
			setAccountHeaders(req, tt.accountId, tt.accountName)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestMovieHandlers_LikeMovie(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		accountId   int
		accountName string
		movieId     string
		liked       interface{}
		wantStatus  int
	}{
		{
			name: "positive_like_movie",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			liked:       false,
			wantStatus:  http.StatusOK,
		},
		{
			name: "positive_delete_like_movie",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			liked:       true,
			wantStatus:  http.StatusOK,
		},
		{
			name: "negative_like_movie_invalid_movie_id",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     invalidId,
			liked:       false,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_like_movie_invalid_liked_param_err",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			liked:       "notBool",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_like_movie_storage_error",
			fields: fields{
				storage: &mock.Storage{
					LikeMovieErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			liked:       false,
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "negative_delete_like_movie_storage_error",
			fields: fields{
				storage: &mock.Storage{
					DeleteMovieLikeErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			movieId:     validId,
			liked:       true,
			wantStatus:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
			)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/movies/%s/like?liked=%v", tt.movieId, tt.liked)
			req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)
			setAccountHeaders(req, tt.accountId, tt.accountName)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestMovieHandlers_ListMovies(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		accountId   int
		accountName string
		query      string
		wantStatus  int
	}{
		{
			name: "positive_list_movies",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			accountId:   1,
			accountName: accountName,
			wantStatus:  http.StatusOK,
		},
		{
			name: "positive_list_movies_without_account_info",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "positive_list_movies_not_default_params",
			fields: fields{
				storage: &mock.Storage{},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			query: "?offset=10&title=mad&order_by=popularity",
			wantStatus:  http.StatusOK,
		},
		{
			name: "negative_list_movies_storage_error",
			fields: fields{
				storage: &mock.Storage{
					ListMoviesErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "negative_list_movies_invalid_params_error",
			fields: fields{
				storage: &mock.Storage{},
				conf:   &config.Config{},
				logger: logger,
			},
			accountId:   1,
			accountName: accountName,
			query: "?offset=invalidOffset",
			wantStatus:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)
			a := NewApi(
				WithConfig(tt.fields.conf),
				WithLogger(tt.fields.logger),
				WithStorage(tt.fields.storage),
			)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/movies%s", tt.query)
			req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)
			setAccountHeaders(req, tt.accountId, tt.accountName)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func setAccountHeaders(r *http.Request, accountId int, accountName string) {
	r.Header.Set(accountHeaderId, fmt.Sprintf("%d", accountId))
	r.Header.Set(accountHeader, accountName)
}

func checkResponseStatusCode(t *testing.T, want int, got int) {
	if want != got {
		t.Errorf("Expected response status code: %d, got: %d", want, got)
	}
}
