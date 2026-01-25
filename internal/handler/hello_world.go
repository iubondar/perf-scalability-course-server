package handler

import "net/http"

type HelloWorldHandler struct{}

func NewHelloWorldHandler() HelloWorldHandler {
	return HelloWorldHandler{}
}

func (h HelloWorldHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
