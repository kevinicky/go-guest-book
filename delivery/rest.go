package delivery

import (
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"net/http"
)

type HTTPHandler struct{}

const apiVersion = "/api/v1"

func (h *HTTPHandler) NewRest(r *mux.Router, healthAdapter adapter.HealthAdapter, userAdapter adapter.UserAdapter, visitAdapter adapter.VisitAdapter, threadAdapter adapter.ThreadAdapter, authAdapter adapter.AuthAdapter) {
	s := r.PathPrefix(apiVersion).Subrouter()

	withMiddleRouter := r.PathPrefix(apiVersion).Subrouter()
	withMiddleRouter.Use(jwtAndMatrixAuth(authAdapter))

	healthRouter(s, healthAdapter)
	userRouter(withMiddleRouter, userAdapter)
	visitRouter(withMiddleRouter, visitAdapter)
	threadRouter(withMiddleRouter, threadAdapter)
	authRouter(s, authAdapter, userAdapter)
}

func healthRouter(s *mux.Router, healthAdapter adapter.HealthAdapter) {
	s.HandleFunc("/health", health(healthAdapter))
}

func userRouter(s *mux.Router, userAdapter adapter.UserAdapter) {
	u := s.PathPrefix("/users").Subrouter()
	u.HandleFunc("/list", getUsers(userAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{user_id}", getUser(userAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{user_id}/delete", deleteUser(userAdapter)).Methods(http.MethodDelete)
	u.HandleFunc("/{user_id}/edit", updateUser(userAdapter)).Methods(http.MethodPut)
}

func visitRouter(s *mux.Router, visitAdapter adapter.VisitAdapter) {
	u := s.PathPrefix("/visits").Subrouter()
	u.HandleFunc("", createVisit(visitAdapter)).Methods(http.MethodPost)
	u.HandleFunc("/list", getVisits(visitAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{visit_id}", getVisit(visitAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{visit_id}/delete", deleteVisit(visitAdapter)).Methods(http.MethodDelete)
}

func threadRouter(s *mux.Router, threadAdapter adapter.ThreadAdapter) {
	u := s.PathPrefix("/threads").Subrouter()
	u.HandleFunc("", createThread(threadAdapter)).Methods(http.MethodPost)
	u.HandleFunc("/list", getThreads(threadAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{thread_id}", getThread(threadAdapter)).Methods(http.MethodGet)
	u.HandleFunc("/{thread_id}/delete", deleteThread(threadAdapter)).Methods(http.MethodDelete)
}

func authRouter(s *mux.Router, authAdapter adapter.AuthAdapter, userAdapter adapter.UserAdapter) {
	u := s.PathPrefix("/auth").Subrouter()
	u.HandleFunc("/login", login(authAdapter)).Methods(http.MethodPost)
	u.HandleFunc("/register", createUser(userAdapter)).Methods(http.MethodPost)
}
