package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &Handler{}
	r := gin.Default()
	r.GET("/ping", h.Ping)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}
	body := decodeJSON[map[string]string](t, w)
	if body["message"] != "pong" {
		t.Fatal("got: ", body["message"], "want: pong")
	}
}
