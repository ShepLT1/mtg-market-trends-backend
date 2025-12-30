package data

import (
	"time"
	"fmt"

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
	PriceDiffPct 			float64 `json:"price_diff_pct"`
	Name        			string  `json:"name"`
	CollectorNum  		string  `json:"collector_num"`
	Finish						string	`json:"finish"`
	ImageURI					string	`json:"image_uri"`
}

type CardWithPrice struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	SetID				 int		 `json:"set_id"`
	ImageURI     string  `json:"image_uri"`
	CollectorNum string  `json:"collector_num"`
	Finish       string  `json:"finish"`
	PromoType    string  `json:"promo_type"`
	AltStyle		 string	 `json:"alt_style"`
	CurrPrice    float64 `json:"curr_price"`
}

// GetListingPriceDiffs fetches listings for a given start/end date and returns top N by price difference
func (r *ListingRepository) GetListingPriceDiffs(startDate, endDate time.Time, limit int, order string) ([]ListingWithDiff, error) {
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	results := make([]ListingWithDiff, limit)

	query := `
		WITH start_min AS (
			SELECT
				c.name,
				MIN(l.price) AS start_price
			FROM listings l
			JOIN cards c ON l.card_id = c.id
			WHERE l.created_date = ?
				AND c.finish = 'nonfoil'
				AND c.alt_style = ''
			GROUP BY c.name
		),
		current_min AS (
			SELECT
				c.name,
				MIN(l.price) AS curr_price
			FROM listings l
			JOIN cards c ON l.card_id = c.id
			WHERE l.created_date = ?
				AND c.finish = 'nonfoil'
				AND c.alt_style = ''
			GROUP BY c.name
		),
		current_cheapest_printing AS (
			SELECT DISTINCT ON (c.name)
				l.id AS listing_id,
				l.card_id,
				c.name,
				l.price,
				c.collector_num,
				c.finish,
				c.image_uri
			FROM listings l
			JOIN cards c ON l.card_id = c.id
			WHERE l.created_date = ?
				AND c.finish = 'nonfoil'
				AND c.alt_style = ''
			ORDER BY c.name, l.price ASC, l.created_at DESC
		)
		SELECT
			cp.listing_id,
			cp.card_id,
			cp.name,
			ROUND(cm.curr_price, 2) AS curr_price,
			ROUND(sm.start_price, 2) AS start_price,
			ROUND(
				((cm.curr_price - sm.start_price) / NULLIF(sm.start_price, 0)) * 100,
				0
			) AS price_diff_pct,
			cp.collector_num,
			cp.finish,
			cp.image_uri
		FROM start_min sm
		JOIN current_min cm
			ON sm.name = cm.name
		JOIN current_cheapest_printing cp
			ON cp.name = sm.name
		WHERE (sm.start_price >= 0.5 OR cm.curr_price >= 0.5)
		ORDER BY price_diff_pct ` + order + `
		LIMIT ?;
	`

	err := r.DB.Raw(
		query,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		limit,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CardRepository) GetCards(name string, page, limit int) ([]CardWithPrice, error) {
	var cards []CardWithPrice
	offset := (page - 1) * limit

	query := `
		SELECT 
			c.id,
			c.name,
			c.set_id,
			c.image_uri,
			c.collector_num,
			c.finish,
			c.promo_type,
			c.alt_style,
			(
				SELECT l.price 
				FROM listings l 
				WHERE l.card_id = c.id 
				ORDER BY l.created_at DESC 
				LIMIT 1
			) AS curr_price
		FROM cards c
		WHERE c.name ILIKE ?
		ORDER BY curr_price DESC
		LIMIT ? OFFSET ?;
	`

	err := r.DB.Raw(query, "%"+name+"%", limit, offset).Scan(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cards: %w", err)
	}

	return cards, nil
}

func (r *CardRepository) GetCardsDistName(name string, limit int) ([]Card, error) {
	var cards []Card

	query := `
		SELECT DISTINCT ON (name)
			id,
			name,
			set_id,
			collector_num,
			promo_type,
			finish,
			alt_style,
			image_uri
		FROM cards
		WHERE name ILIKE ?
		ORDER BY name, id
		LIMIT ?;
	`

	err := r.DB.Raw(query, "%"+name+"%", limit).Scan(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cards: %w", err)
	}

	return cards, nil
}