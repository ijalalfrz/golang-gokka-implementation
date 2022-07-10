package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// CORS returns cors middleware.
func CORS(handler http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}),
	)(handler)
}
