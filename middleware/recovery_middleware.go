package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

// RecoveryLogger is a struct
type RecoveryLogger struct {
	Logger *logrus.Logger
}

// Println will log the recovered panic
func (rl *RecoveryLogger) Println(recovered ...interface{}) {
	rl.Logger.Error(recovered...)
}

// Recovery is a recovery middleware.
func Recovery(logger *logrus.Logger, handler http.Handler) http.Handler {
	rl := &RecoveryLogger{Logger: logger}
	return handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(true),
		handlers.RecoveryLogger(rl),
	)(handler)
}
