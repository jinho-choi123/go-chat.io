package hubManager

import (
	"context"
	"log"
	"strconv"

	"go-chat.io/src/redisManager"
	"go-chat.io/src/safeWebsocket"
)

type Hub struct {
	HubID      uint64
	Publish    chan []byte
	Broadcast  chan []byte
	Clients    map[*safeWebsocket.SafeWebsocket]bool
	Register   chan *safeWebsocket.SafeWebsocket
	Unregister chan *safeWebsocket.SafeWebsocket
}

func NewHub(hubId uint64) *Hub {
	return &Hub{
		HubID:      hubId,
		Publish:    make(chan []byte),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*safeWebsocket.SafeWebsocket]bool),
		Register:   make(chan *safeWebsocket.SafeWebsocket),
		Unregister: make(chan *safeWebsocket.SafeWebsocket),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true

		case client := <-hub.Unregister:
			log.Printf("Hub(hubId=%v) Unregistering client...", hub.HubID)
			if _, ok := hub.Clients[client]; ok {
				delete(hub.Clients, client)
				close(client.Send)
			}

		case message := <-hub.Broadcast:
			for client := range hub.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(hub.Clients, client)
				}
			}
		}
	}
}

// act as redis publisher thread for the hub
func (hub *Hub) RedisPub(rdm *redisManager.RedisManager) {
	// return connection to the pool when function ends
	conn := rdm.Pool.Conn()
	defer conn.Close()

	ctx := context.Background()
	channel := strconv.FormatUint(hub.HubID, 10)

	for {
		select {
		case message := <-hub.Publish:
			conn.Publish(ctx, channel, message)
		}
	}
}

func (hub *Hub) RedisSub(rdm *redisManager.RedisManager) {
	rdb := rdm.Pool

	ctx := context.Background()
	channel := strconv.FormatUint(hub.HubID, 10)

	subscriber := rdb.Subscribe(ctx, channel)
	defer subscriber.Close()

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			log.Panic(err.Error())
		}
		// double check if message channel is right
		if msg.Channel != channel {
			log.Panic("message channel is incorrect")
		}
		hub.Broadcast <- []byte(msg.Payload)
	}
}
