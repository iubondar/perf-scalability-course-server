package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHelloWorldHandler(t *testing.T) {
	h := NewHelloWorldHandler()
	assert.NotNil(t, h)
	assert.IsType(t, HelloWorldHandler{}, h)
}

func TestHelloWorldHandler_Handle(t *testing.T) {
	h := NewHelloWorldHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Hello")
	assert.Contains(t, w.Body.String(), "World")
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>")
}

