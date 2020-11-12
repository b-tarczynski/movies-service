package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarTar213/movies-service/models"
	"github.com/gin-gonic/gin"
)

func TestCheckAccount(t *testing.T) {
	tests := []struct {
		name       string
		account    *models.AccountInfo
		wantStatus int
	}{
		{
			name: "positive_check_account",
			account: &models.AccountInfo{
				ID:    1,
				Login: "login",
				Role:  "role",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "positive_check_account_invalid_headers",
			account: &models.AccountInfo{
				ID:    1,
				Role:  "role",
			},
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/ping", nil)
			setAccountHeaders(req, tt.account)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CheckAccount() = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CheckAccount())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}

func setAccountHeaders(r *http.Request, account *models.AccountInfo) {
	if account == nil {
		return
	}

	r.Header.Set("X-Account-Id", fmt.Sprintf("%d", account.ID))
	r.Header.Set("X-Account", account.Login)
	r.Header.Set("X-Role", account.Role)
}
