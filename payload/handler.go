package payload

import (
	"fmt"
	"ws_pubsub/pubsub"
)

type Payload struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
	Client  pubsub.Client
}

type WsResposne struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

const SUBSCRIBE = "subscribe"
const UNSUBSCRIBE = "unsubscribe"
const PUBLISH = "publish"

func HandlePayload(client pubsub.Client, c chan<- Payload) {
	payload := Payload{}
	for {
		err := client.Connection.ReadJSON(&payload)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("There is payload: %v\n", payload)
			payload.Client = client
			c <- payload
		}
	}
}

func HandleWsChannel(c <-chan Payload, ps *pubsub.Pubsub) {
	for {
		payload := <-c
		switch payload.Action {
		case SUBSCRIBE:
			if payload.Topic == "" {
				fmt.Println("Cannot Subscribe!!. There is no topice")
				break
			}
			subscription := pubsub.Subscription{
				Topic:  payload.Topic,
				Client: payload.Client,
			}
			fmt.Printf("topic: %v\n", subscription.Topic)
			ps.Subscriptions = append(ps.Subscriptions, subscription)
		case UNSUBSCRIBE:
			fmt.Println("unsubscibe")
			if payload.Topic == "" {
				fmt.Println("Cannot Unsubscribe!!. There is no topice")
				break
			}
			unsubscibe(payload.Topic, payload.Client, ps)
		case PUBLISH:
			fmt.Println("publish")
			if payload.Topic == "" || payload.Message == "" {
				fmt.Println("There is no topic or no message")
				break
			}
			res := WsResposne{
				Topic:   payload.Topic,
				Message: payload.Message,
			}
			publishToTopic(res, payload.Client, ps)
		}
	}
}

func publishToTopic(res WsResposne, client pubsub.Client, ps *pubsub.Pubsub) {
	for _, sub := range ps.Subscriptions {
		if sub.Topic == res.Topic && sub.Client.Id != client.Id {
			err := sub.Client.Connection.WriteJSON(res)
			if err != nil {
				fmt.Printf("Send message error: %v\n", err)
			}
		}
	}
}

func unsubscibe(topic string, client pubsub.Client, ps *pubsub.Pubsub) {
	for index, sub := range ps.Subscriptions {
		if sub.Topic == topic && sub.Client.Id == client.Id {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
}
