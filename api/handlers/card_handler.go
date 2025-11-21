package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/data"
)

// GetPriceDiffsHandler returns a handler function with the repository injected
func GetCards(cardRepo *data.CardRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		page := 1
		limit := 50
		if p := r.URL.Query().Get("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				page = parsed
			}
		}
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 500 {
				limit = parsed
			}
		}
		name := r.URL.Query().Get("name")

		cards, err := cardRepo.GetCards(name, page, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"page":  page,
			"limit": limit,
			"data":  cards,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
