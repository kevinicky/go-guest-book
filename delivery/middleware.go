package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"net/http"
	"strings"
)

func jwtAndMatrixAuth(authAdapter adapter.AuthAdapter) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			jwtToken := r.Header.Get("Authorization")
			respJwt, err := authAdapter.ValidateJWT(jwtToken)
			if err != nil {
				jsonResp, _ := json.Marshal(entity.ErrorMessage{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				})
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(jsonResp)

				return
			}

			userMatrixReq := entity.UserMatrixValidateRequest{
				UserID:   respJwt.Credential,
				Endpoint: r.URL.String(),
			}

			for _, muxVar := range mux.Vars(r) {
				userMatrixReq.Endpoint = strings.ReplaceAll(userMatrixReq.Endpoint, muxVar, "[id]")
			}
			userMatrixReq.Endpoint = strings.Split(userMatrixReq.Endpoint, "?")[0]

			err = authAdapter.ValidateUserMatrixAccess(userMatrixReq)
			if err != nil {
				var statusCode int
				switch err.Error() {
				case customerror.UNAUTHORISED_ACCESS:
					statusCode = http.StatusUnauthorized
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

			h.ServeHTTP(w, r)
		})
	}
}
