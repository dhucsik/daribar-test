package utils

import (
	"testing"

	"github.com/dhucsik/daribar-test/model"
)

func TestNewFanout(t *testing.T) {
	f := NewFanout()

	if f.channels == nil {
		t.Errorf("Channels map was not initialized")
	}

	if f.connected {
		t.Errorf("Expected connected to be false, got true")
	}
}

func TestConnected(t *testing.T) {
	f := &Fanout{
		channels:  make(map[string]model.Client),
		connected: false,
	}

	f.Connected()

	if !f.connected {
		t.Errorf("Expected connected to be true, got false")
	}
}

func TestDisconnected(t *testing.T) {
	f := &Fanout{
		channels: map[string]model.Client{
			"123": {Ch: make(chan int), Quantity: 1},
			"456": {Ch: make(chan int), Quantity: 2},
		},
		connected: true,
	}

	f.Disconnected()

	if f.connected {
		t.Errorf("Expected connected to be false, got true")
	}

	if len(f.channels) != 0 {
		t.Errorf("Expected channels to be empty, got %v", f.channels)
	}
}

func TestSubscribe(t *testing.T) {
	f := &Fanout{
		channels:  make(map[string]model.Client),
		connected: true,
	}

	ch := make(chan int)
	f.Subscribe("123", 1, ch)

	if len(f.channels) != 1 {
		t.Errorf("Expected one channel, got %v", len(f.channels))
	}

	client, ok := f.channels["123"]
	if !ok {
		t.Errorf("Expected channel for phone number 123, got none")
	}

	if client.Quantity != 1 {
		t.Errorf("Expected quantity to be 1, got %v", client.Quantity)
	}

	if client.Ch != ch {
		t.Errorf("Channels do not match")
	}
}

func TestSubscribeNotConnected(t *testing.T) {
	f := &Fanout{
		channels:  make(map[string]model.Client),
		connected: false,
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got none")
		}
	}()

	f.Subscribe("123", 1, make(chan int))
}

func TestUnsubscribe(t *testing.T) {
	f := &Fanout{
		channels: map[string]model.Client{
			"123": {Ch: make(chan int), Quantity: 1},
			"456": {Ch: make(chan int), Quantity: 2},
		},
		connected: true,
	}

	f.Unsubscribe("123")

	if len(f.channels) != 1 {
		t.Errorf("Expected one channel, got %v", len(f.channels))
	}

	_, ok := f.channels["123"]
	if ok {
		t.Errorf("Expected no channel for phone number 123, got one")
	}
}

func TestUnsubscribeNonExistent(t *testing.T) {
	f := &Fanout{
		channels:  make(map[string]model.Client),
		connected: true,
	}

	f.Unsubscribe("123")
}
