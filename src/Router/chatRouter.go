package router

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat.io/src/hubManager"
	"go-chat.io/src/safeWebsocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ChatEndpoint(hubTable *hubManager.HubTable) func(*gin.Context) {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Error occur at /ws: %v", err)
			return
		}

		hubid_query, ok := c.Params.Get("hubId")
		if !ok {
			log.Printf("Incorrect HubID %v", hubid_query)
			return
		}

		hubID, err := strconv.ParseUint(hubid_query, 10, 64)
		if err != nil {
			log.Printf("Incorrect HubID %v", hubid_query)
			return
		}

		target_hub, err := hubTable.Lookup(hubID)
		if err != nil {
			log.Printf("HubTable Lookup failed: %v", err)
			return
		}

		sws := safeWebsocket.NewSafeWebsocket(conn)
		target_hub.Register <- sws
		go sws.Read(target_hub.Publish, target_hub.Unregister)
		go sws.Write()
	}
}
