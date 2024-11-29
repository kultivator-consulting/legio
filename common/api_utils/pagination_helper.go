package api_utils

import (
	"cortex_api/common/api_models"
	"github.com/gofiber/fiber/v2"
	"math"
	"strconv"
)

func GetPaginationParams(ctx *fiber.Ctx) (int, int, string, string) {
	query, err := ParseQueryString(ctx)
	if err != nil {
		return 1, 10, "created", "asc"
	}

	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("itemsPerPage"))
	sortBy := query.Get("sortBy")
	sortOrder := query.Get("sortOrder")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	return page, pageSize, sortBy, sortOrder
}

func GetPaginatedResults(totalRows int64, requestedPage int, requestedPageSize int, data interface{}) api_models.Results {
	totalPages := int(math.Ceil(float64(totalRows) / float64(requestedPageSize)))

	previousPage := requestedPage - 1
	if previousPage <= 0 {
		previousPage = 1
	}
	nextPage := requestedPage + 1
	if nextPage >= totalPages {
		nextPage = totalPages
	}

	return api_models.Results{
		Status:       ResultsResponse,
		TotalRows:    totalRows,
		TotalPages:   totalPages,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		Data:         data,
	}
}
