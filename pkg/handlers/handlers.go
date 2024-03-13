package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var clients = make(map[*websocket.Conn]UserData)
var broadcast = make(chan WsPayload)

type UserData struct {
	CurrentApp string `json:"current_app"`
	Username   string `json:"username"`
}

type WsJsonResponse struct {
	Action         string     `json:"action"`
	Message        string     `json:"message"`
	MessageType    string     `json:"message_type"`
	ConnectedUsers []UserData `json:"connected_users"`
}

type WsPayload struct {
	Action     string          `json:"action"`
	Username   string          `json:"username"`
	Message    string          `json:"message"`
	CurrentApp string          `json:"current_app"`
	Conn       *websocket.Conn `json:"-"`
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
			clients[e.Conn] = UserData{
				Username:   e.Username,
				CurrentApp: "",
			}
			users := getUserList()

			response.Action = "list_users"
			response.Message = e.Message
			response.ConnectedUsers = users

			log.Printf("users list: %+v\n", users)

			broadcastToAll(&response)
		case "left":
			response.Action = "left"
			delete(clients, e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(&response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("%s : %s", e.Username, e.Message)
			broadcastToAll(&response)
		case "appInfo":
			response.Action = "app_info"
			response.Message = fmt.Sprintf("%s,%s", e.Username, e.CurrentApp)
			broadcastToAll(&response)
		}
	}
}

func broadcastToAll(response *WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Error broadcasting JSON:", err)
			client.WriteMessage(websocket.CloseMessage, []byte{})
			client.Close()
		}
	}
}

func getUserList() []UserData {
	var userList []UserData
	for _, x := range clients {
		if x.Username != "" {
			userList = append(userList, x)
		}
	}
	return userList
}

func WsEndpoint(c *websocket.Conn) {
	defer func() {
		c.Close()
	}()

	clients[c] = UserData{
		Username:   "",
		CurrentApp: "",
	}
	log.Println("Client connected to endpoint")

	var payload WsPayload

	for {
		// err := c.ReadJSON(&payload)
		messageType, p, err := c.ReadMessage()
		if err != nil {
			// do nothing
			break
		}
		if messageType == websocket.TextMessage {
			err := json.Unmarshal(p, &payload)
			if err != nil {
				//
				break
			}
			payload.Conn = c
			broadcast <- payload
		}
	}
}
