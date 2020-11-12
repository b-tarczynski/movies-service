package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/BarTar213/movies-service/config"
	"github.com/BarTar213/movies-service/mock"
	"github.com/BarTar213/movies-service/models"
	"github.com/BarTar213/movies-service/storage"
	"github.com/gin-gonic/gin"
)

const (
	contentExample = "example content"
)

func TestNewCommentHandlers(t *testing.T) {
	type args struct {
		storage storage.Storage
		logger  *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *CommentHandlers
	}{
		{
			name: "positiveNewCommentHandlers",
			args: args{
				storage: &mock.Storage{},
				logger:  &log.Logger{},
			},
			want: &CommentHandlers{
				storage: &mock.Storage{},
				logger:  &log.Logger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommentHandlers(tt.args.storage, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentHandlers_AddComment(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		account     *models.AccountInfo
		movieId     string
		body        interface{}
		wantStatus  int
	}{
		{
			name: "positive_add_comment",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			movieId:     validId,
			body: &models.Comment{
				Content: contentExample,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "negative_add_comment_invalid_movie_id_param_error",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			movieId:     invalidId,
			body: &models.Comment{
				Content: contentExample,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_add_comment_storage_error",
			fields: fields{
				storage: &mock.Storage{
					AddMovieCommentErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			movieId:     validId,
			body: &models.Comment{
				Content: contentExample,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "negative_add_comment_invalid_body",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			movieId:     validId,
			body: map[string]interface{}{
				"content": 12,
			},
			wantStatus: http.StatusBadRequest,
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

			jsonBody, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/comments?movie_id=%s", tt.movieId)
			req, _ := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(jsonBody))
			setAccountHeaders(req, tt.account)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCommentHandlers_DeleteComment(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		account    *models.AccountInfo
		commentId  string
		wantStatus int
	}{
		{
			name: "positive_delete_comment",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			wantStatus:  http.StatusOK,
		},
		{
			name: "negative_delete_comment_invalid_comment_id_error",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   invalidId,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_delete_comment_storage_error",
			fields: fields{
				storage: &mock.Storage{
					DeleteCommentErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
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
			reqUrl := fmt.Sprintf("/comments/%s", tt.commentId)
			req, _ := http.NewRequest(http.MethodDelete, reqUrl, nil)
			setAccountHeaders(req, tt.account)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCommentHandlers_GetComments(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		account    *models.AccountInfo
		query       string
		movieId     string
		wantStatus  int
	}{
		{
			name: "positive_get_comments",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			movieId:     validId,
			wantStatus:  http.StatusOK,
		},
		{
			name: "positive_get_comments_not_default_pagination_params",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			query:       "&order_by=likes&offset=2&limit=20",
			movieId:     validId,
			wantStatus:  http.StatusOK,
		},
		{
			name: "positive_get_comments_invalid_pagination_params",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			query:       "&order_by=likes&offset=2&limit=invalidLimit",
			movieId:     validId,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_get_comments_invalid_movie_id_param_error",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			movieId:     invalidId,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_get_comments_storage_error",
			fields: fields{
				storage: &mock.Storage{
					GetMovieCommentsErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
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
			reqUrl := fmt.Sprintf("/comments?movie_id=%s%s", tt.movieId, tt.query)
			req, _ := http.NewRequest(http.MethodGet, reqUrl, nil)
			setAccountHeaders(req, tt.account)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCommentHandlers_LikeComment(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		account    *models.AccountInfo
		commentId   string
		liked       interface{}
		wantStatus  int
	}{
		{
			name: "positive_like_comment",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			liked:       false,
			wantStatus:  http.StatusOK,
		},
		{
			name: "positive_delete_like_comment",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			liked:       true,
			wantStatus:  http.StatusOK,
		},
		{
			name: "negative_like_comment_invalid_comment_id",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   invalidId,
			liked:       false,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_like_comment_invalid_liked_param_err",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			liked:       "notBool",
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "negative_like_comment_storage_error",
			fields: fields{
				storage: &mock.Storage{
					LikeCommentErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			liked:       false,
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name: "negative_delete_like_comment_storage_error",
			fields: fields{
				storage: &mock.Storage{
					DeleteCommentLikeErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
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
			reqUrl := fmt.Sprintf("/comments/%s/like?liked=%v", tt.commentId, tt.liked)
			req, _ := http.NewRequest(http.MethodPost, reqUrl, nil)
			setAccountHeaders(req, tt.account)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func TestCommentHandlers_UpdateComment(t *testing.T) {
	type fields struct {
		storage storage.Storage
		conf    *config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		account    *models.AccountInfo
		commentId   string
		body        interface{}
		wantStatus  int
	}{
		{
			name: "positive_update_comment",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			body: &models.Comment{
				Content: contentExample,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "negative_update_comment_invalid_comment_id",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   invalidId,
			body: &models.Comment{
				Content: contentExample,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_update_comment_invalid_body_error",
			fields: fields{
				storage: &mock.Storage{},
				conf:    &config.Config{},
				logger:  logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			body: map[string]interface{}{
				"content": 12,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative_update_comment_storage_error",
			fields: fields{
				storage: &mock.Storage{
					UpdateCommentErr: true,
				},
				conf:   &config.Config{},
				logger: logger,
			},
			account: &models.AccountInfo{
				ID:    1,
				Login: accountLogin,
				Role:  accountRole,
			},
			commentId:   validId,
			body: &models.Comment{
				Content: contentExample,
			},
			wantStatus: http.StatusInternalServerError,
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

			jsonBody, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			reqUrl := fmt.Sprintf("/comments/%s", tt.commentId)
			req, _ := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(jsonBody))
			setAccountHeaders(req, tt.account)

			a.Router.ServeHTTP(w, req)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}
