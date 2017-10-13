package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/di0nys1us/httpgo"
)

var (
	h *Handler
)

func TestMain(m *testing.M) {
	h = &Handler{}
	os.Exit(m.Run())
}

func TestGetUsers(t *testing.T) {
	handler := httpgo.ErrorHandlerFunc(h.GetUsers)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(w, r)
}
