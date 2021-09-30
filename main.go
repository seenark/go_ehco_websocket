package main

import (
	"fmt"
	"log"
	"net/http"
	"ws_pubsub/payload"
	"ws_pubsub/pubsub"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var ps = pubsub.Pubsub{
	Clients: []pubsub.Client{},
}

var WsChannel = make(chan payload.Payload)

func main() {
	e := echo.New()
	e.Static("/public", "./public")
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Hello": "World",
		})
	})

	e.GET("/ws", websocketHandler)

	go payload.HandleWsChannel(WsChannel, &ps)

	if err := e.Start("127.0.0.1:8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func websocketHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	fmt.Println("Client connected")

	client := pubsub.Client{
		Id:         uuid.NewString(),
		Connection: ws,
	}

	ps.Clients = append(ps.Clients, client)
	fmt.Printf("clients: %v\n", ps.Clients)
	payload.HandlePayload(client, WsChannel)
	return nil
}
