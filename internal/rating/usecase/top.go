package usecase

import (
	"time"

	"github.com/zenoleg/binomeme/internal/rating"
)

const (
	TopDay TopPreset = iota
	TopWeek
	TopMonth
	TopEver
)

type (
	TopPreset uint8

	TopMemesQuery struct {
		now    time.Time
		period TopPreset
	}

	TopMemesHandler struct {
		storage rating.MemeStorage
	}
)

func NewTopMemesQuery(now time.Time, period TopPreset) TopMemesQuery {
	return TopMemesQuery{
		now:    now,
		period: period,
	}
}
