package middleware

import (
	"beetle/internal/domain"
	"beetle/internal/handler"
	"fmt"
	"net/http"

	"strconv"

	"github.com/labstack/echo/v4"
)

func Paginate(hc *handler.Context) error {
	hc.PaginationQuery = &domain.PaginationQuery{
		Page:  ParseIntQueryParam(hc, "page", 1),
		Limit: ParseIntQueryParam(hc, "limit", domain.PageLimitDefault),
	}
	if hc.PaginationQuery.Limit > domain.PageLimitMax || hc.PaginationQuery.Limit < domain.PageLimitMin {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("limit must be between 1 and %d", domain.PageLimitMax))
	}
	if hc.PaginationQuery.Page < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "page must be greater than 0")
	}
	return nil
}

func ParseIntQueryParam(hc *handler.Context, name string, fallback int) int {
	strParam := hc.Context.QueryParam(name)
	if strParam != "" {
		converted, err := strconv.Atoi(strParam)
		// On error, the fallback will be returned
		if err == nil {
			return converted
		}
	}
	return fallback
}
