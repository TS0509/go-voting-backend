package middleware

import (
	"net"
	"net/http"

	"sync"

	"golang.org/x/time/rate"
)

// IPLimiter 存储每个 IP 的限速器
type IPLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		visitors: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (i *IPLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.visitors[ip] = limiter
	}
	return limiter
}

// RateLimitMiddleware 速率限制中间件
func (i *IPLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		limiter := i.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
