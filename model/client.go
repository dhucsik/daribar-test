package model

type Client struct {
	Ch       chan<- int
	Quantity int
}
