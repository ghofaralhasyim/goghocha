package handlers

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan WsPayload)

type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string          `json:"action"`
	Username string          `json:"username"`
	Message  string          `json:"message"`
	Conn     *websocket.Conn `json:"-"`
}

func Homepage(c *fiber.Ctx) error {
	tmpl := c.Render("index", fiber.Map{}, "layouts/base")

	return tmpl
}

type WebSocketConnection struct {
	*websocket.Conn
}

func ListenToWs() {
	var response WsJsonResponse
	for {
		e := <-broadcast
		switch e.Action {
		case "username":
			clients[e.Conn] = e.Username
			users := getUserList()

			response.Action = "list_users"
			response.Message = e.Message
			response.ConnectedUsers = users

			log.Printf("users list: %+v\n", users)

			broadcastToAll(response)
		case "left":
			response.Action = "left"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<span class='user'>%s :</span><span>%s</span>", e.Username, e.Message)
			broadcastToAll(response)
		}
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Error broadcasting JSON:", err)
			client.WriteMessage(websocket.CloseMessage, []byte{})
			client.Close()
		}
	}
}

func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	return userList
}

func WsEndpoint(c *websocket.Conn) {
	defer func() {
		c.Close()
	}()

	clients[c] = ""
	log.Println("Client connected to endpoint")

	var payload WsPayload

	for {
		err := c.ReadJSON(&payload)
		if err != nil {
			// do nothing
			break
		} else {
			payload.Conn = c
			broadcast <- payload
		}
	}
}
