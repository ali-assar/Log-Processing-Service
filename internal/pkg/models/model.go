package models

import "errors"

type Order struct {
	ID        string   `json:"id"`
	Amount    float64  `json:"amount"`
	Items     []string `json:"items"`
	Customer  string   `json:"customer"`
	Status    string   `json:"status"`
	CreatedAt int64    `json:"created_at"`
	Address   string   `json:"address"`
	Notes     string   `json:"notes,omitempty"`
}

var validStatuses = map[string]bool{
	"pending":   true,
	"paid":      true,
	"shipped":   true,
	"delivered": true,
	"cancelled": true,
}

// Validate checks whether the order has all required fields with acceptable values.
func (o *Order) Validate() error {
	if o.Customer == "" {
		return errors.New("customer is required")
	}
	if o.Status == "" {
		return errors.New("status is required")
	}
	if !validStatuses[o.Status] {
		return errors.New("invalid status")
	}
	if o.Address == "" {
		return errors.New("address is required")
	}
	if o.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if len(o.Items) == 0 {
		return errors.New("items must not be empty")
	}
	return nil
}
