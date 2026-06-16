package handler

import (
	"druna_server/pkg/repository"
	"druna_server/pkg/service"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests-only")
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestPingEndpoint(t *testing.T) {
	services := service.NewService(&repository.Repository{})
	h := NewHandler(services)
	router := h.InitRoutes()

	req := httptest.NewRequest(http.MethodGet, "/ping/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
