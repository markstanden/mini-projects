package http

import (
	"net/http"

	"github.com/markstanden/authentication"
)

// Our wrapper of the Handler interface
type Handler struct {
	UserService authentication.UserService
}

// Implements the Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle request
}

func (h *Handler) ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
