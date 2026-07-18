package domain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusCreated   Status = "created"
	StatusPaid      Status = "paid"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

type Order struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Name      string    `json:"name"`
	Total     float64   `json:"total"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var ErrOrderNotFound = errors.New("order not found")
