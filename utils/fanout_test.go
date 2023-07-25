package utils

import (
	"fmt"
	"testing"
	"time"

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

/*
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
}*/

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

func TestSubscribe(t *testing.T) {
	f := &Fanout{
		channels:  make(map[string]model.Client),
		connected: true,
	}

	ch := make(chan int)
	for i := 0; i < 100000; i++ {
		f.Subscribe(fmt.Sprintf("%v", i), 1, ch)
	}

	for _, m := range f.channels {
		m.Ch <- 1
	}
	client, ok := f.channels["1234"]
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

func TestPublish(t *testing.T) {
	// Initialize Fanout struct
	f := &Fanout{
		channels:  make(map[string]model.Client),
		connected: true,
	}

	// Scenario 1: Test with update == nil
	ch := make(chan int)
	f.channels["1234567890"] = model.Client{
		Ch:       ch,
		Quantity: 1,
	}
	go f.Publish(nil)
	select {
	case <-time.After(time.Second):
		t.Fatal("Expected to receive quantity, but didn't")
	case q, ok := <-ch:
		if !ok {
			t.Errorf("Channel closed unexpectedly")
		}
		if q != 1 {
			t.Errorf("Expected Quantity 1, got %d", q)
		}
	}

	// Scenario 2: Test with a valid update and Increment
	phone := "123"
	ch = make(chan int)
	f.channels[phone] = model.Client{Ch: ch, Quantity: 1}

	update := &model.Update{
		Order: model.Order{Phone: phone},
		Inc:   true,
	}

	go f.Publish(update)

	select {
	case <-time.After(time.Second):
		t.Fatal("Expected to receive update, but didn't")
	case quantity := <-ch:
		if quantity != 2 {
			t.Errorf("Expected quantity to be 2, got %v", quantity)
		}
	}

	// Scenario 3: Test with a valid update and Decrement
	phone = "1122334455"
	ch = make(chan int)
	f.channels[phone] = model.Client{
		Ch:       ch,
		Quantity: 2,
	}
	update = &model.Update{
		Order: model.Order{
			Phone: phone,
		},
		Inc: false,
	}
	go f.Publish(update)
	select {
	case <-time.After(time.Second):
		t.Fatal("Expected to receive update, but didn't")
	case q, ok := <-ch:
		if !ok {
			t.Errorf("Channel closed unexpectedly")
		}
		if q != 1 {
			t.Errorf("Expected Quantity 1, got %d", q)
		}
	}

	// Scenario 4: Test with an invalid phone number in update (causing a panic)
	phone = "0000000000"
	update = &model.Update{
		Order: model.Order{
			Phone: phone,
		},
		Inc: false,
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()
	f.Publish(update)
}
