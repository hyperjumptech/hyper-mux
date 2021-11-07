package hyper_mux

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCORSMiddleware(t *testing.T) {
	mux := NewHyperMux()
	mux.UseMiddleware(NewCORSMiddleware(DefaultCORSOption))
	mux.AddRoute("/testing", MethodGet, func(writer http.ResponseWriter, request *http.Request) {
		WriteString(writer, http.StatusOK, "OK")
	})

	recorder := httptest.NewRecorder()
	requst, err := http.NewRequest(MethodOptions, "http://serv/testing", nil)
	requst.Header.Add("Origin", "https://other.com")
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, requst)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	requst, err = http.NewRequest(MethodGet, "http://serv/testing", nil)
	requst.Header.Add("Origin", "https://other.com")
	assert.NoError(t, err)
	mux.ServeHTTP(recorder, requst)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "OK", recorder.Body.String())
}
