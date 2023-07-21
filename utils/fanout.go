package utils

import (
	"errors"
	"log"
	"sync"

	"github.com/dhucsik/daribar-test/model"
)

type Fanout struct {
	m         sync.Mutex
	channels  map[string]model.Client
	connected bool
}

func NewFanout() *Fanout {
	return &Fanout{channels: make(map[string]model.Client)}
}

func (f *Fanout) Connected() {
	f.m.Lock()
	f.connected = true
	f.m.Unlock()
	log.Println("connected")
}

func (f *Fanout) Disconnected() {
	f.m.Lock()
	f.connected = false
	for c := range f.channels {
		close(f.channels[c].Ch)
		delete(f.channels, c)
	}
	f.m.Unlock()
	log.Println("disconnected")
}

func (f *Fanout) Subscribe(phone string, quantity int, c chan<- int) {
	f.m.Lock()
	defer f.m.Unlock()
	if !f.connected {
		panic(errors.New("fanout: not connected"))
	}

	f.channels[phone] = model.Client{
		Ch:       c,
		Quantity: quantity,
	}
	log.Println("Subscribed")
}

func (f *Fanout) Unsubscribe(phone string) {
	f.m.Lock()
	_, ok := f.channels[phone]
	if ok {
		close(f.channels[phone].Ch)
		delete(f.channels, phone)
	}
	f.m.Unlock()
	log.Println("unSubscribed")
}

func (f *Fanout) Publish(update *model.Update) {
	f.m.Lock()
	if update == nil {
		for c := range f.channels {
			select {
			case f.channels[c].Ch <- f.channels[c].Quantity:
			default:
				close(f.channels[c].Ch)
				delete(f.channels, c)
			}
		}
	} else {
		client, ok := f.channels[update.Order.Phone]
		if !ok {
			panic(errors.New("fanout not connected"))
		}

		if update.Inc {
			client.Quantity = client.Quantity + 1
		} else {
			client.Quantity = client.Quantity - 1
		}

		f.channels[update.Order.Phone] = client
		select {
		case f.channels[update.Order.Phone].Ch <- f.channels[update.Order.Phone].Quantity:
		default:
			close(f.channels[update.Order.Phone].Ch)
			delete(f.channels, update.Order.Phone)
		}
	}
	f.m.Unlock()
}
