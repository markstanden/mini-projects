package http

import (
	"net/http"

	"github.com/markstanden/authentication/authentication"
)

type Handler struct {
        UserService authentication.UserService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        // handle request
}