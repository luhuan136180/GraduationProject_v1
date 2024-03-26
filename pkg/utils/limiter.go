package utils

import (
	"golang.org/x/time/rate"
	"time"
)

var LmimiterMap = make(map[string]*rate.Limiter)
var Lmt = NewLimiter(rate.Every(time.Second), 5)

type limiter struct {
	limit rate.Limit
	burst int
}

func NewLimiter(r rate.Limit, b int) *limiter {
	return &limiter{
		limit: r,
		burst: b,
	}
}

func (l *limiter) AllowKey(key string) bool {
	limiter, ok := LmimiterMap[key]
	if !ok {
		// create limiter
		limiter = rate.NewLimiter(l.limit, l.burst)
		LmimiterMap[key] = limiter
	}

	return limiter.Allow()
}
