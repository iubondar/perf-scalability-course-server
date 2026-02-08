package handlers

import (
	_ "embed"
	"net/http"
)

//go:embed templates/hello.html
var helloHTML string

type HelloWorldHandler struct{}

func NewHelloWorldHandler() HelloWorldHandler {
	return HelloWorldHandler{}
}

func (h HelloWorldHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(helloHTML))
}
