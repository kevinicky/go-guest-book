package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kevinicky/go-guest-book/internal/adapter"
	"github.com/kevinicky/go-guest-book/internal/entity"
	"github.com/kevinicky/go-guest-book/util"
	"github.com/kevinicky/go-guest-book/util/customerror"
	"net/http"
	"strings"
)

func createVisit(visitAdapter adapter.VisitAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		var payload entity.CreateVisitRequest
		if err := decoder.Decode(&payload); err != nil {
			msg := map[string]string{"error": err.Error()}
			jsonResp, _ := json.Marshal(msg)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonResp)

			return
		}

		resp, errList := visitAdapter.CreateVisit(payload)

		if len(errList) > 0 {
			var httpStatusCode int
			var errResp []entity.ErrorMessage

			for _, err := range errList {
				switch err.Error() {
				case customerror.USER_ID_IS_MANDATORY:
					httpStatusCode = http.StatusBadRequest
				case customerror.USER_NOT_FOUND:
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

func getVisit(visitAdapter adapter.VisitAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		visitID := mux.Vars(r)["visit_id"]
		resp, err := visitAdapter.GetVisit(visitID)

		if err != nil {
			var httpStatusCode int

			if strings.HasPrefix(err.Error(), "uuid:") {
				httpStatusCode = http.StatusBadRequest
			} else {
				switch err.Error() {
				case customerror.VISIT_NOT_FOUND:
					httpStatusCode = http.StatusNotFound
				default:
					httpStatusCode = http.StatusInternalServerError
				}
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

func getVisits(visitAdapter adapter.VisitAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		limit, offset := util.GetLimitOffset(*r)
		resp, err := visitAdapter.GetVisits(limit, offset)

		if err != nil {
			httpStatusCode := http.StatusInternalServerError

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

func deleteVisit(visitAdapter adapter.VisitAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		visitID := mux.Vars(r)["visit_id"]
		err := visitAdapter.DeleteVisit(visitID)

		if err != nil {
			var httpStatusCode int

			if strings.HasPrefix(err.Error(), "uuid:") {
				httpStatusCode = http.StatusBadRequest
			} else if err.Error() == customerror.VISIT_NOT_FOUND {
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

		resp := map[string]string{
			"message": customerror.VISIT_SUCCESS_DELETE,
		}

		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResp)
	}
}
