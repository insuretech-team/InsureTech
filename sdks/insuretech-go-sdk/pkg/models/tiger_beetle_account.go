package models

import (
	"time"
)

// TigerBeetleAccount represents a tiger_beetle_account
type TigerBeetleAccount struct {
	AccountId string `json:"account_id"`
	AccountType string `json:"account_type"`
	Balance string `json:"balance"`
	BalanceUpdatedAt time.Time `json:"balance_updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Currency string `json:"currency"`
	IsActive bool `json:"is_active"`
	LedgerId int `json:"ledger_id"`
	TigerbeetleAccountId string `json:"tigerbeetle_account_id"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId string `json:"user_id,omitempty"`
}
