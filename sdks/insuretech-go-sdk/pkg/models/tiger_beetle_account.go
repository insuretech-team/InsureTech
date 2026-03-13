package models

import (
	"time"
)

// TigerBeetleAccount represents a tiger_beetle_account
type TigerBeetleAccount struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TigerbeetleAccountId string `json:"tigerbeetle_account_id"`
	AccountType string `json:"account_type"`
	LedgerId int `json:"ledger_id"`
	Currency string `json:"currency"`
	BalanceUpdatedAt time.Time `json:"balance_updated_at,omitempty"`
	AccountId string `json:"account_id"`
	UserId string `json:"user_id,omitempty"`
	Balance string `json:"balance"`
	IsActive bool `json:"is_active"`
}
