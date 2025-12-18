package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend/data"
)

// GetPriceDiffsHandler returns a handler function with the repository injected
func GetPriceDiffsHandler(listingRepo *data.ListingRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Default values
		limit := 10
		order := "desc"

		// Default 1-day period
		endDate := time.Now().UTC()
		startDate := endDate.AddDate(0, 0, -1)

		// Parse query parameters
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		if o := r.URL.Query().Get("order"); o == "asc" || o == "desc" {
			order = o
		}
		if s := r.URL.Query().Get("start_date"); s != "" {
			if parsed, err := time.Parse("2006-01-02", s); err == nil {
				startDate = parsed
			}
		}
		if e := r.URL.Query().Get("end_date"); e != "" {
			if parsed, err := time.Parse("2006-01-02", e); err == nil {
				endDate = parsed
			}
		}

		results, err := listingRepo.GetListingPriceDiffs(startDate, endDate, limit, order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)

	}
}
