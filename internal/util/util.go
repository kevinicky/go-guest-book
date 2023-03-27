package util

import (
	"math"
	"net/http"
	"strconv"
)

func CountTotalPageAndCurrentPage(totalRows int64, limit, offset int) (int64, int64) {
	var totalPages int64
	var page int64

	totalPages = int64(math.Ceil(float64(totalRows) / float64(limit)))
	page = int64(math.Floor(float64(offset)/float64(limit))) + 1

	return totalPages, page
}

func GetLimitOffset(r http.Request) (int, int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	return limit, offset
}
