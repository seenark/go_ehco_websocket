package pubsub

import "github.com/gorilla/websocket"

type Client struct {
	Id         string
	Connection *websocket.Conn
}

type Pubsub struct {
	Clients       []Client
	Subscriptions []Subscription
}

type Subscription struct {
	Topic  string
	Client Client
}
