package delivery

import (
	"encoding/json"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"net/http"
)

func createUser(userAdapter adapter.UserAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		decoder := json.NewDecoder(r.Body)
		var payload entity.CreateUserRequest
		if err := decoder.Decode(&payload); err != nil {
			msg := map[string]string{"error": err.Error()}
			jsonResp, _ := json.Marshal(msg)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonResp)

			return
		}

		resp, errList := userAdapter.CreateUser(payload)

		if len(errList) > 0 {
			var httpStatusCode int
			var errResp []entity.ErrorMessage

			for _, err := range errList {
				switch err.Error() {
				case "password cannot be empty":
					httpStatusCode = http.StatusBadRequest
				case "password cannot more than 64 characters":
					httpStatusCode = http.StatusBadRequest
				case "email is not valid":
					httpStatusCode = http.StatusBadRequest
				case "email has taken":
					httpStatusCode = http.StatusBadRequest
				case "phone number has taken":
					httpStatusCode = http.StatusBadRequest
				case "full_name is mandatory":
					httpStatusCode = http.StatusBadRequest
				case "phone_number is mandatory":
					httpStatusCode = http.StatusBadRequest
				default:
					httpStatusCode = http.StatusInternalServerError
				}

				errResp = append(errResp, entity.ErrorMessage{
					Code:    httpStatusCode,
					Message: err.Error(),
				})
			}

			msg := map[string][]entity.ErrorMessage{"error": errResp}
			jsonResp, _ := json.Marshal(msg)
			w.WriteHeader(httpStatusCode)
			_, _ = w.Write(jsonResp)

			return
		}

		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResp)
	}
}
