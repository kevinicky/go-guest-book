package delivery

import (
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"net/http"
)

type HTTPHandler struct{}

const apiVersion = "/api/v1"

func (h *HTTPHandler) NewRest(r *mux.Router, healthAdapter adapter.HealthAdapter, userAdapter adapter.UserAdapter) {
	s := r.PathPrefix(apiVersion).Subrouter()

	healthRouter(s, healthAdapter)
	userRouter(s, userAdapter)
}

func healthRouter(s *mux.Router, healthAdapter adapter.HealthAdapter) {
	s.HandleFunc("/health", health(healthAdapter))
}

func userRouter(s *mux.Router, userAdapter adapter.UserAdapter) {
	u := s.PathPrefix("/users").Subrouter()
	u.HandleFunc("", createUser(userAdapter)).Methods(http.MethodPost)
	u.HandleFunc("/list", getUsers(userAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{user_id}", getUser(userAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{user_id}/delete", deleteUser(userAdapter)).Methods(http.MethodDelete)
	u.HandleFunc("/{user_id}/edit", updateUser(userAdapter)).Methods(http.MethodPut)
}
