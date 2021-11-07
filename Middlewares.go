package hyper_mux

import (
	"bytes"
	"context"
	"github.com/rs/cors"
	"math/rand"
	"net/http"
	"time"
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

type ContextKey string

const (
	RequestID ContextKey = "REQUEST-ID"
)

// ContextSetterMiddleware is a middleware function that will ensure the request context existance.
func ContextSetterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx := context.WithValue(context.Background(), RequestID, MakeRequestID())
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetRequestID will retrieve the RequestID value from context if ContextSetterMiddleware middleware is used.
// If ContextSetterMiddleware is not used, it will return empty string
func GetRequestID(r *http.Request) string {
	ctx := r.Context()
	if ctx == nil {
		return ""
	}
	iv := ctx.Value(RequestID)
	if iv == nil {
		return ""
	}
	if s, ok := iv.(string); ok {
		return s
	}
	return ""
}

// NewCORSMiddleware will create CORS middleware to be used in your web app
// this middleware will handle the OPTIONS method automatically.
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

const (
	CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func MakeRequestID() string {
	rand.Seed(time.Now().UnixMicro())
	buff := &bytes.Buffer{}
	for buff.Len() < 20 {
		offset := rand.Intn(len(CharSet))
		buff.WriteString(CharSet[offset : offset+1])
	}
	return buff.String()
}
