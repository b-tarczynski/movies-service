package api

import (
	"net/http"
	"strconv"

	"github.com/BarTar213/movies-service/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

func (h *MovieHandlers) GetCredits(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(movieIdKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	credits := &models.Credit{}
	err = h.storage.GetCredits(id, credits)
	if err != nil && err != pg.ErrNoRows {
		handlePostgresError(c, h.logger, err, creditsResource)
		return
	}

	if err == pg.ErrNoRows {
		status, err := h.tmdb.GetCredits(id, credits)
		if err != nil || status != http.StatusOK {
			handleTMDBError(c, h.logger, status, err, creditsResource)
			return
		}
		credits.MovieId = id
		go h.AddCredits(credits)
	}

	c.JSON(http.StatusOK, credits)
}

func (h *MovieHandlers) AddCredits(credits *models.Credit) {
	for i := 0; i < len(credits.Crew); i++ {
		credits.Crew[i].CreditId = credits.Id
	}
	for i := 0; i < len(credits.Cast); i++ {
		credits.Cast[i].CreditId = credits.Id
	}

	err := h.storage.AddCredits(credits)
	if err != nil {
		err := errors.WithMessage(err, "AddCredits")
		h.logger.Print(err)
	}
}
