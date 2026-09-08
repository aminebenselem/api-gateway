package ratelimiter

import "golang.org/x/time/rate"

type RateLimiter struct {
	limiter *rate.Limiter
}

func New(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}