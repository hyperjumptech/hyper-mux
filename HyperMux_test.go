package hyper_mux

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHyperMux_UseMiddleware(t *testing.T) {
	mux := NewHyperMux()
	steps := 0
	mux.AddRoute("/ahoy", "GET", func(w http.ResponseWriter, r *http.Request) {
		t.Log("AHOY ENDPOINT")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("AHOY"))
		steps++
	})

	mux.UseMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("AHOY MIDDLEWARE 1")
			steps++
			next.ServeHTTP(w, r)
		})
	})

	mux.UseMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("AHOY MIDDLEWARE 2")
			steps++
			next.ServeHTTP(w, r)
		})
	})

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://dummy/ahoy", nil)
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, 3, steps)
}

func TestHyperMux_UseMiddleware2(t *testing.T) {
	mux := NewHyperMux()
	steps := 0
	mux.AddRoute("/ahoy", "GET", func(w http.ResponseWriter, r *http.Request) {
		t.Log("AHOY ENDPOINT")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("AHOY"))
		steps++
	})

	mux.UseMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("AHOY MIDDLEWARE 1")
			steps++
		})
	})

	mux.UseMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("AHOY MIDDLEWARE 2")
			next.ServeHTTP(w, r)
			steps++
		})
	})

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://dummy/ahoy", nil)
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, 1, steps)
}

func TestHyperMux_Pattern(t *testing.T) {
	mux := NewHyperMux()
	mux.AddRoute("/ahoy/{somekey}/hoya", "GET", func(w http.ResponseWriter, r *http.Request) {
		t.Log("AHOY ENDPOINT")

		assert.Equal(t, "somevalue", r.Header.Get("somekey"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("AHOY"))
	})

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://dummy/ahoy/somevalue/hoya", nil)
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHyperMux_ServeHTTP(t *testing.T) {

}
