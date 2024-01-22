package entity

import "time"

type Account struct {
	Id     int
	Wallet []Wallet `json:"wallet"`
}

type Wallet struct {
	Id            int
	ActiveBalance float64 `json:"active_balance"`
	FrozenBalance float64 `json:"frozen_balance"`
	Currency      string  `json:"currency"`
}

type Transaction struct {
	Id         int
	WalletId   int       `json:"wallet_id,omitempty"`
	Status     string    `json:"status,omitempty"`
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	Withdraw   bool      `json:"withdraw,omitempty"`
	CardNumber string    `json:"card_number,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// User request a transaction
type RequestTransaction struct {
	UserId int         `json:"user_id"`
	Trans  Transaction `json:"transaction"`
}
