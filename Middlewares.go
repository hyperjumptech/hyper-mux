package hyper_mux

import (
	"context"
	"github.com/rs/cors"
	"net/http"
)

var (
	DefaultCORSOption = cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{MethodPost, MethodGet, MethodDelete, MethodPut},
		AllowedHeaders:     []string{"Authorization", "Content-Type", "Content-Length", "Content-Encoding", "Accept", "Accept-Encoding"},
		ExposedHeaders:     []string{"*", "Authorization"},
		MaxAge:             300,
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              true,
	}
	Cors *cors.Cors
)

// ContextSetterMiddleware is a middleware function that will ensure the request context existance.
func ContextSetterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			next.ServeHTTP(w, r.WithContext(context.Background()))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NewCORSMiddleware will create CORS middleware to be used in your web app
func NewCORSMiddleware(options cors.Options) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if Cors == nil {
			Cors = cors.New(options)
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == MethodOptions {
				ohandler := Cors.Handler(&corsOptionHandler{})
				ohandler.ServeHTTP(w, r)
				return
			}
			chandler := Cors.Handler(next)
			chandler.ServeHTTP(w, r)
		})
	}
}

type corsOptionHandler struct {
}

func (h *corsOptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
