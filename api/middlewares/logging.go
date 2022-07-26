package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	nested "github.com/algo-matchfund/grants-backend/internal/nested_logger"
)

func NewLoggingMiddleware(stackTraceEnabled bool, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var logBuffer strings.Builder
		var logger nested.NestedTraceLogger
		logger.SetOutput(&logBuffer)
		logger.SetFlags(log.LstdFlags | log.Lmsgprefix)
		logger.SetPrefix(fmt.Sprintf("%s %s: ", r.Method, r.URL.EscapedPath()))
		handlerStart := time.Now()

		// Dump accumulated log in case of panic
		defer func(start time.Time, buffer *strings.Builder) {
			err := recover()
			if err != nil {
				http.Error(rw, "Internal server error", http.StatusInternalServerError)
				logger.Println("panicked")
			}

			// the log message might have been created with Println, we need to append duration on the same string
			// so remove end of line character if it's found
			logger.Println("duration: " + time.Since(start).String())
			logString := buffer.String()
			if stackTraceEnabled && err != nil {
				logString += "\n" + string(debug.Stack())
			}
			fmt.Print(logString)
		}(handlerStart, &logBuffer)

		handler.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), "logger", logger)))
	})
}
