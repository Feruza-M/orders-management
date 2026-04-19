package models

import "time"

type Order struct {
	ID          int64     `json:"id"`
	Customer    string    `json:"customer"`
	Email       string    `json:"email"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	Customer    string  `json:"customer" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Amount      float64 `json:"amount" binding:"required"`
	Status      string  `json:"status" binding:"required"`
	Description string  `json:"description"`
}
