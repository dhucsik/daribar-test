package model

import "time"

type Order struct {
	OrderID   int
	Phone     string
	CreatedAt time.Time
	IsOpen    bool
}
