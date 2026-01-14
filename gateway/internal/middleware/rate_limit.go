package middleware

import (
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	limiter  *rate.Limiter
	lastSeen int64 // Nanosecond timestamp
}

type RateLimiterStats struct {
	TotalRequests   atomic.Int64
	BlockedRequests atomic.Int64
	CacheSize       atomic.Int64
}

// RateLimiter middleware: IP-based throttling (using LRU cache)
//
// Parameters:
//   - rps: Requests per second (token generation rate)
//   - burst: Maximum burst requests (bucket capacity)
//
// Features:
//   - Use LRU cache to automatically evict inactive IPs
//   - Support X-Forwarded-For、X-Real-IP proxy headers
//   - Record lastSeen for statistics and debugging
func RateLimiter(rps int, burst int) func(next http.Handler) http.Handler {
	// LRU cache, max 10000 IPs
	cache, err := lru.New[string, *ClientLimiter](10_000)
	if err != nil {
		panic("failed to create LRU cache: " + err.Error())
	}

	stats := &RateLimiterStats{}

	// Get or create limiter
	getLimiter := func(ip string) *ClientLimiter {
		now := time.Now().UnixNano()

		// Get from cache
		if cl, ok := cache.Get(ip); ok {
			// Update last seen time
			atomic.StoreInt64(&cl.lastSeen, now)
			return cl
		}

		// Create new limiter
		cl := &ClientLimiter{
			limiter:  rate.NewLimiter(rate.Limit(rps), burst),
			lastSeen: now,
		}
		cache.Add(ip, cl)
		stats.CacheSize.Store(int64(cache.Len()))

		return cl
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			stats.TotalRequests.Add(1)

			// Get client IP
			ip := getClientIP(r)

			// Get limiter
			cl := getLimiter(ip)

			// Check if request is allowed
			if !cl.limiter.Allow() {
				stats.BlockedRequests.Add(1)

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-RateLimit-Limit", string(rune(rps)))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Rate limit exceeded","code":429}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP Extracts the real client IP from the request
//
// Priority:
//  1. X-Real-IP (nginx/cloudflare setting)
//  2. X-Forwarded-For's first IP (standard proxy header)
//  3. RemoteAddr (direct connection)
func getClientIP(r *http.Request) string {
	// 1. Use X-Real-IP if available
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	// 2. Parse X-Forwarded-For (multiple IPs possible, format: client, proxy1, proxy2)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Get first IP (client's real IP)
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 3. Use RemoteAddr (may contain port)
	ip := r.RemoteAddr
	// Remove port: "192.168.1.1:12345" -> "192.168.1.1"
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}
