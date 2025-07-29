package middleware

import (
	"go-web/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiterMiddleware creates a rate limiter middleware.
func RateLimiterMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Parse the rate limit period from the configuration.
	period, err := time.ParseDuration(cfg.RateLimiter.Period)
	if err != nil {
		panic("Invalid rate limiter period")
	}

	// Create a new rate limiter with a memory store.
	rate := limiter.Rate{
		Period: period,
		Limit:  cfg.RateLimiter.Limit,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get the IP address of the client.
		ip := c.ClientIP()
		// Create a context for the rate limiter.
		context, err := instance.Get(c, ip)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Check if the request is allowed.
		if context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		// If the request is allowed, call the next handler.
		c.Next()
	}
}
