package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseRecorder) Write(p []byte) (int, error) {
	// Ensure status is set even if WriteHeader wasn't explicitly called
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(p)
	rw.bytes += n
	return n, err
}

func clientIP(r *http.Request) string {
	// Prefer X-Forwarded-For (first IP), then X-Real-IP, then RemoteAddr.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		dur := time.Since(start)
		log.Printf("%s %s %d %dB %s ip=%s ua=%q",
			r.Method,
			r.URL.RequestURI(),
			rec.status,
			rec.bytes,
			dur.Round(time.Millisecond),
			clientIP(r),
			r.UserAgent(),
		)
	})
}
