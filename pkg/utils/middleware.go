package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// A custom response writer to capture the response
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	startTime  time.Time
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// Calculate elapsed time and set the header when Write is called
	elapsedTime := time.Since(rw.startTime)
	rw.Header().Set("X-Request-Time", strconv.FormatInt(int64(elapsedTime), 10))
	return rw.ResponseWriter.Write(b)
}

func TimeOfReqRespCycle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		timeOfReq := time.Now()
		fmt.Println(timeOfReq)
		resWriter := &responseWriter{
			writer,
			http.StatusOK,
			timeOfReq,
		}
		next.ServeHTTP(resWriter, request)

	})
}

func BearerTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")

		if authHeader != "" && !strings.HasPrefix(authHeader, "Bearer ") {
			request.Header.Set("Authorization", "Bearer "+authHeader)
		}
		next.ServeHTTP(writer, request)
	})
}
