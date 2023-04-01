package delivery

import (
	"encoding/json"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"net/http"
)

func login(authAdapter adapter.AuthAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		credential, password, ok := r.BasicAuth()
		if !ok {
			jsonResp, _ := json.Marshal(entity.ErrorMessage{
				Code:    http.StatusBadRequest,
				Message: customerror.AUTH_NOT_SUPPLIED,
			})
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonResp)

			return
		}

		err := authAdapter.CheckCredentials(credential, password)
		if err != nil {
			var statusCode int
			if err.Error() == customerror.INVALID_CREDENTIAL {
				statusCode = http.StatusBadRequest
			} else {
				statusCode = http.StatusInternalServerError
			}

			jsonResp, _ := json.Marshal(entity.ErrorMessage{
				Code:    statusCode,
				Message: err.Error(),
			})
			w.WriteHeader(statusCode)
			_, _ = w.Write(jsonResp)

			return
		}

		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/json" {
			jsonResp, _ := json.Marshal(entity.ErrorMessage{
				Code:    http.StatusBadRequest,
				Message: customerror.INVALID_JSON_HEADER,
			})
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonResp)

			return
		}

		var payload entity.JwtRequest
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&payload); err != nil {
			jsonResp, _ := json.Marshal(entity.ErrorMessage{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonResp)

			return
		}

		payload.Credential = credential

		jwtToken, err := authAdapter.CreateJWT(payload)
		if err != nil {
			var statusCode int
			switch err.Error() {
			case customerror.INVALID_CREDENTIAL:
				statusCode = http.StatusBadRequest
			case customerror.ISSUER_MANDATORY:
				statusCode = http.StatusBadRequest
			default:
				statusCode = http.StatusInternalServerError
			}

			jsonResp, _ := json.Marshal(entity.ErrorMessage{
				Code:    statusCode,
				Message: err.Error(),
			})
			w.WriteHeader(statusCode)
			_, _ = w.Write(jsonResp)

			return
		}

		jsonResp, _ := json.Marshal(jwtToken)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResp)
	}
}
