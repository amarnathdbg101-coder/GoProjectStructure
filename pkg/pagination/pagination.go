package pagination

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Params contains extracted query parameters for paginated queries.
type Params struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`  // "asc" or "desc"
	Search string `json:"search"` // sanitized search keyword
}

// Meta contains metadata returned in paginated API responses.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Envelope is the standard wrapper for paginated lists.
type Envelope struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
	Meta    Meta `json:"meta"`
}

// GetParams parses query parameters from an HTTP request with safe defaults.
// Default: page=1, limit=10 (max limit 100), order=desc, sort_by=created_at.
func GetParams(r *http.Request) Params {
	q := r.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100 // Prevent fetching excessive rows
	}

	offset := (page - 1) * limit

	order := strings.ToLower(strings.TrimSpace(q.Get("order")))
	if order != "asc" {
		order = "desc"
	}

	sortBy := strings.TrimSpace(q.Get("sort_by"))
	if sortBy == "" {
		sortBy = "created_at"
	}

	search := strings.TrimSpace(q.Get("search"))

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
		Search: search,
	}
}

// BuildMeta calculates pagination metadata based on total row count.
func BuildMeta(totalItems int64, params Params) Meta {
	totalPages := 0
	if params.Limit > 0 {
		totalPages = int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return Meta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}

// Respond sends a standardized 200 OK JSON response with pagination metadata.
func Respond(w http.ResponseWriter, items any, totalItems int64, params Params) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	meta := BuildMeta(totalItems, params)

	_ = json.NewEncoder(w).Encode(Envelope{
		Success: true,
		Data:    items,
		Meta:    meta,
	})
}
