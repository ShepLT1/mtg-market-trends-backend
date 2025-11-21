package data

import (
	"time"

	"gorm.io/gorm"
)

// CardRepository defines methods for querying cards
type CardRepository struct {
	DB *gorm.DB
}

// NewCardRepository creates a new CardRepository
func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{DB: db}
}

type ListingRepository struct {
	DB *gorm.DB
}

// NewListingRepository creates a new ListingRepository
func NewListingRepository(db *gorm.DB) *ListingRepository {
	return &ListingRepository{DB: db}
}

// ListingWithDiff represents a listing and its price difference
type ListingWithDiff struct {
	ListingID       	int64   `json:"listing_id"`
	CardID          	int64   `json:"card_id"`
	CurrPrice        	float64 `json:"curr_price"`
	StartPrice      	float64 `json:"start_price"`
	PriceDiff 				float64 `json:"price_diff"`
	Name        			string  `json:"name"`
	CollectorNum  		string  `json:"collector_num"`
	Finish						string	`json:"finish"`
	ImageURI					string	`json:"image_uri"`
}

// GetListingPriceDiffs fetches listings for a given start/end date and returns top N by price difference
func (r *ListingRepository) GetListingPriceDiffs(startDate, endDate time.Time, limit int, order string) ([]ListingWithDiff, error) {
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	var results []ListingWithDiff

	query := `
	SELECT 
		curr.listing_id,
		curr.card_id,
		curr.price AS curr_price,
		start.price AS start_price,
		(curr.price - start.price) AS price_diff,
		c.name,
		c.collector_num,
		c.finish,
		c.image_uri
	FROM 
		(SELECT DISTINCT ON (card_id) id AS listing_id, card_id, price FROM listings WHERE created_date = ? ORDER BY card_id, created_at DESC
		) curr
	JOIN 
		(SELECT DISTINCT ON (card_id) card_id, price FROM listings WHERE created_date = ? ORDER BY card_id, created_at ASC
		) start
	ON curr.card_id = start.card_id
	JOIN cards c ON curr.card_id = c.id
	WHERE c.finish = 'nonfoil'
  AND c.alt_style IS NULL
	ORDER BY price_diff ` + order + `
	LIMIT ?;
	`

	err := r.DB.Raw(
		query,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		limit,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetCards returns cards optionally filtered by exact name, with pagination
func (r *CardRepository) GetCards(name string, page, limit int) ([]Card, error) {
	var cards []Card
	offset := (page - 1) * limit

	dbQuery := r.DB.Model(&Card{})

	if name != "" {
		dbQuery = dbQuery.Where("name = ?", name)
	}

	result := dbQuery.Limit(limit).Offset(offset).Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}

	return cards, nil
}
