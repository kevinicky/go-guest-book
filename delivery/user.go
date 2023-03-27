package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/customerror"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/internal/util"
	"net/http"
	"strings"
)

func createUser(userAdapter adapter.UserAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(jsonResp)
	}
}

func getUser(userAdapter adapter.UserAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := mux.Vars(r)["user_id"]
		resp, err := userAdapter.GetUser(userID)

		if err != nil {
			var httpStatusCode int

			if strings.HasPrefix(err.Error(), "uuid:") {
				httpStatusCode = http.StatusBadRequest
			} else if err.Error() == customerror.USER_NOT_FOUND {
				httpStatusCode = http.StatusNotFound
			} else {
				httpStatusCode = http.StatusInternalServerError
			}

			msg := entity.ErrorMessage{
				Code:    httpStatusCode,
				Message: err.Error(),
			}

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

func getUsers(userAdapter adapter.UserAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		limit, offset := util.GetLimitOffset(*r)
		key := r.URL.Query().Get("key")
		isAdmin := r.URL.Query().Get("is_admin")

		resp, err := userAdapter.GetUsers(limit, offset, key, isAdmin)

		if err != nil {
			var httpStatusCode int
			if err.Error() == customerror.IS_ADMIN_WRONG_CONTENT {
				httpStatusCode = http.StatusBadRequest
			} else {
				httpStatusCode = http.StatusInternalServerError
			}

			msg := entity.ErrorMessage{
				Code:    httpStatusCode,
				Message: err.Error(),
			}

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
