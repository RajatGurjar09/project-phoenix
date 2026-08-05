package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var requestLogger = log.New(os.Stdout, "", 0)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(body)
}

// RequestLogger logs one structured JSON record for every HTTP request.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		responseWriter := &loggingResponseWriter{ResponseWriter: w}

		defer func() {
			statusCode := responseWriter.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			entry := map[string]any{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"method":    r.Method,
				"path":      r.URL.Path,
				"status":    statusCode,
				"duration":  time.Since(startedAt).String(),
				"client_ip": clientIP(r),
			}

			payload, err := json.Marshal(entry)
			if err == nil {
				requestLogger.Println(string(payload))
			}
		}()

		next.ServeHTTP(responseWriter, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
