package middleware

import (
	"net/http"

	"github.com/didip/tollbooth/v7"

	"github.com/artarts36/specw"
)

type RateLimitConfig struct {
	Max specw.Env[float64] `yaml:"max" json:"max"`
}

func (c *RateLimitConfig) Validate() error {
	const defaultRateLimit = 100

	if c.Max.Value == 0 {
		c.Max.Value = defaultRateLimit
	}

	return nil
}

func RateLimit(next http.Handler, config RateLimitConfig) http.Handler {
	tlbthLimiter := tollbooth.NewLimiter(config.Max.Value, nil)

	return tollbooth.LimitFuncHandler(tlbthLimiter, next.ServeHTTP)
}
