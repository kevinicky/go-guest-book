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

func createThread(threadAdapter adapter.ThreadAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		var payload entity.CreateThreadRequest
		if err := decoder.Decode(&payload); err != nil {
			msg := map[string]string{"error": err.Error()}
			jsonResp, _ := json.Marshal(msg)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(jsonResp)

			return
		}

		resp, errList := threadAdapter.CreateThread(payload)

		if len(errList) > 0 {
			var httpStatusCode int
			var errResp []entity.ErrorMessage

			for _, err := range errList {
				switch err.Error() {
				case customerror.VISIT_ID_IS_MANDATORY:
					httpStatusCode = http.StatusBadRequest
				case customerror.VISIT_NOT_FOUND:
					httpStatusCode = http.StatusBadRequest
				case customerror.CONTENT_IS_MANDATORY:
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

func getThread(threadAdapter adapter.ThreadAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		threadID := mux.Vars(r)["thread_id"]
		resp, err := threadAdapter.GetThread(threadID)

		if err != nil {
			var httpStatusCode int

			if strings.HasPrefix(err.Error(), "uuid:") {
				httpStatusCode = http.StatusBadRequest
			} else {
				switch err.Error() {
				case customerror.THREAD_NOT_FOUND:
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

func getThreads(threadAdapter adapter.ThreadAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		limit, offset := util.GetLimitOffset(*r)
		visitID := r.URL.Query().Get("visit_id")
		if visitID == "" {
			httpStatusCode := http.StatusBadRequest

			msg := entity.ErrorMessage{
				Code:    httpStatusCode,
				Message: customerror.VISIT_ID_IS_MANDATORY,
			}

			jsonResp, _ := json.Marshal(msg)
			w.WriteHeader(httpStatusCode)
			_, _ = w.Write(jsonResp)

			return
		}

		resp, err := threadAdapter.GetThreads(limit, offset, visitID)

		if err != nil {
			var httpStatusCode int
			if strings.HasPrefix(err.Error(), "uuid:") {
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

func deleteThread(threadAdapter adapter.ThreadAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		threadID := mux.Vars(r)["thread_id"]
		err := threadAdapter.DeleteThread(threadID)

		if err != nil {
			var httpStatusCode int

			if strings.HasPrefix(err.Error(), "uuid:") {
				httpStatusCode = http.StatusBadRequest
			} else if err.Error() == customerror.THREAD_NOT_FOUND {
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
			"message": customerror.THREAD_SUCCESS_DELETE,
		}

		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(jsonResp)
	}
}
