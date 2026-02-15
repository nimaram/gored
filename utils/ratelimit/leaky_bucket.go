package ratelimit

import (
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

type bucket struct {
	tokens int
	last   time.Time
}

var (
	mutex      sync.Mutex
	buckets = map[string]*bucket{}
)

const (
	capacity = 10
	leakRate = time.Second
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mutex.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &bucket{tokens: capacity, last: time.Now()}
			buckets[ip] = b
		}

		elapsed := time.Since(b.last)
		if elapsed >= leakRate {
			b.tokens = capacity
			b.last = time.Now()
		}

		if b.tokens <= 0 {
			mutex.Unlock()
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		b.tokens--
		mutex.Unlock()

		c.Next()
	}
}
