package delivery

import (
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"net/http"
)

type HTTPHandler struct{}

const apiVersion = "/api/v1"

func (h *HTTPHandler) NewRest(r *mux.Router, guestBookAdapter adapter.GuestBookAdapter, userAdapter adapter.UserAdapter) {
	s := r.PathPrefix(apiVersion).Subrouter()
	s.HandleFunc("/health", health(guestBookAdapter))
	user(s, userAdapter)
}

func user(s *mux.Router, userAdapter adapter.UserAdapter) {
	u := s.PathPrefix("/users").Subrouter()
	u.HandleFunc("", createUser(userAdapter)).Methods(http.MethodPost)
	u.HandleFunc("/list", getUsers(userAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{user_id}", getUser(userAdapter)).Methods(http.MethodGet)
}
