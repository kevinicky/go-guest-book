package delivery

import (
	"encoding/json"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/customerror"
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
				case customerror.PASSWORD_LEN_GT_LIMIT:
					httpStatusCode = http.StatusBadRequest
				case customerror.PASSWORD_MANDATORY:
					httpStatusCode = http.StatusBadRequest
				case customerror.INVALID_EMAIL:
					httpStatusCode = http.StatusBadRequest
				case customerror.EMAIL_TAKEN:
					httpStatusCode = http.StatusBadRequest
				case customerror.EMAIL_MANDATORY:
					httpStatusCode = http.StatusBadRequest
				case customerror.PHONE_NUMBER_TAKEN:
					httpStatusCode = http.StatusBadRequest
				case customerror.PHONE_NUMBER_MANDATORY:
					httpStatusCode = http.StatusBadRequest
				case customerror.FULL_NAME_MANDATORY:
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
