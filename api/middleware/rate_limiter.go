package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*Visitor)
	mu       sync.Mutex
)

// RateLimiter membatasi login maksimal 5 kali per menit dari IP yang sama
func RateLimiter() gin.HandlerFunc {
	// Background routine membersihkan log visitor lama
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 1*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &Visitor{lastSeen: time.Now(), count: 1}
			mu.Unlock()
			c.Next()
			return
		}

		v.lastSeen = time.Now()
		v.count++
		if v.count > 5 {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Batas percobaan login terlampaui. Coba lagi dalam 1 menit."})
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}
