package models

import "time"

type Rating struct {
	UserId     int       `json:"user_id" pg:",pk"`
	MovieId    int       `json:"movie_id" pg:",pk"`
	CreateDate time.Time `json:"create_date"`
	Rating     *int      `json:"rating" binding:"required"`
}
