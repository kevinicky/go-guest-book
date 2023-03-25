package delivery

import (
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"net/http"
)

type HTTPHandler struct{}

func (h *HTTPHandler) NewRest(mux *http.ServeMux, guestBookAdapter adapter.GuestBookAdapter) {
	apiVersion := "/api/v1"

	mux.HandleFunc(apiVersion+"/health", health(guestBookAdapter))
}
