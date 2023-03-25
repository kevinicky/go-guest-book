package delivery

import (
	"encoding/json"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"net/http"
)

func health(guestBookAdapter adapter.GuestBookAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		resp := guestBookAdapter.Health()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(jsonResp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}
}
