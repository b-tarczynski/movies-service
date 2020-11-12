package api

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarTar213/movies-service/mock"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

const (
	resourceExample = "resource"
)

func Test_handlePostgresError(t *testing.T) {
	type args struct {
		logger   *log.Logger
		err      error
		resource string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "no_rows_error",
			args: args{
				logger:   logger,
				err:      pg.ErrNoRows,
				resource: resourceExample,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "23503_error",
			args: args{
				logger:   logger,
				err:      &mock.PgError{FieldCode: "23503"},
				resource: resourceExample,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "23505_error",
			args: args{
				logger:   logger,
				err:      &mock.PgError{FieldCode: "23505"},
				resource: resourceExample,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not_handled_postgres_error",
			args: args{
				logger:   logger,
				err:      &mock.PgError{FieldCode: "23111"},
				resource: resourceExample,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)

			handlePostgresError(context, tt.args.logger, tt.args.err, tt.args.resource)

			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}

func Test_handleTMDBError(t *testing.T) {
	type args struct {
		logger   *log.Logger
		status   int
		err      error
		resource string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "negative_handle_tmdb_error_internal_server_error",
			args: args{
				logger:   logger,
				status:   http.StatusInternalServerError,
				err:      errors.New("TMDB error"),
				resource: creditsResource,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "negative_handle_tmdb_status_not_found",
			args: args{
				logger:   logger,
				status:   http.StatusNotFound,
				err:      nil,
				resource: creditsResource,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "negative_handle_tmdb_internal_server_error",
			args: args{
				logger:   logger,
				status:   http.StatusInternalServerError,
				err:      nil,
				resource: creditsResource,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(w)

			handleTMDBError(context, tt.args.logger, tt.args.status, tt.args.err, tt.args.resource)
			checkResponseStatusCode(t, tt.wantStatus, w.Code)
		})
	}
}
