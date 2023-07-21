package storage

import (
	"testing"
)

func TestGetOpenOrdersQuantity(t *testing.T) {
	s := &Storage{}

	quantity := s.GetOpenOrdersQuantity("123")

	if quantity != 10 {
		t.Errorf("Expected quantity to be 10, got %v", quantity)
	}
}
