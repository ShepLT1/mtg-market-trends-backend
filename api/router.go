package api

import (
	"net/http"
	"github.com/gorilla/mux"

	"backend/api/handlers"
	"backend/data"
)

func NewRouter(cardRepo *data.CardRepository, listingRepo *data.ListingRepository) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/cards", handlers.GetCards(cardRepo)).Methods("GET")
	r.HandleFunc("/api/listings", handlers.GetPriceDiffsHandler(listingRepo)).Methods("GET")

	return r
}
