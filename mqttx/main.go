package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Room      string    `json:"room"`
	ToUser    string    `json:"toUser,omitempty"` // 私聊目标用户
	ChatType  string    `json:"chatType"`         // "group" 或 "private"
	Source    string    `json:"source,omitempty"`
}

type Client struct {
	conn     *websocket.Conn
	send     chan Message
	username string
	room     string
	id       string
}

type OnlineUser struct {
	Username string `json:"username"`
	Room     string `json:"room"`
	ID       string `json:"id"`
}

type Hub struct {
	clients     map[*Client]bool
	broadcast   chan Message
	register    chan *Client
	unregister  chan *Client
	privateMsg  chan Message
	userList    chan *Client
	mutex       sync.RWMutex
	onlineUsers map[string]*Client // username -> client
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	hub        *Hub
	mqttClient mqtt.Client
)

func newHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		privateMsg:  make(chan Message),
		userList:    make(chan *Client),
		onlineUsers: make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.onlineUsers[client.username] = client
			h.mutex.Unlock()

			log.Printf("Client %s joined room %s", client.username, client.room)

			// 发送欢迎消息到群聊
			if client.room != "" {
				welcomeMsg := Message{
					Type:      "system",
					Content:   client.username + " joined the chat",
					Timestamp: time.Now(),
					Room:      client.room,
					ChatType:  "group",
				}
				h.broadcastToRoom(welcomeMsg)
			}

			// 广播用户列表更新
			h.broadcastUserList()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.onlineUsers, client.username)
				close(client.send)
				h.mutex.Unlock()

				// 发送离开消息到群聊
				if client.room != "" {
					leaveMsg := Message{
						Type:      "system",
						Content:   client.username + " left the chat",
						Timestamp: time.Now(),
						Room:      client.room,
						ChatType:  "group",
					}
					h.broadcastToRoom(leaveMsg)
				}

				// 广播用户列表更新
				h.broadcastUserList()
			} else {
				h.mutex.Unlock()
			}

		case message := <-h.broadcast:
			// 群聊消息处理
			if message.ChatType == "group" {
				if message.Source == "websocket" {
					message.Source = "mqtt"
					if token := mqttClient.Publish("chat/group/"+message.Room, 0, false, messageToJSON(message)); token.Wait() && token.Error() != nil {
						log.Printf("MQTT publish error: %v", token.Error())
					}
				} else {
					h.broadcastToRoom(message)
				}
			}

		case message := <-h.privateMsg:
			// 私聊消息处理
			if message.ChatType == "private" {
				if message.Source == "websocket" {
					message.Source = "mqtt"
					// 私聊消息发布到特定主题
					topic := "chat/private/" + message.Username + "/" + message.ToUser
					if token := mqttClient.Publish(topic, 0, false, messageToJSON(message)); token.Wait() && token.Error() != nil {
						log.Printf("MQTT publish error: %v", token.Error())
					}
				} else {
					h.sendPrivateMessage(message)
				}
			}

		case client := <-h.userList:
			// 发送用户列表给特定客户端
			h.sendUserListToClient(client)
		}
	}
}

func (h *Hub) broadcastToRoom(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.room == message.Room {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) sendPrivateMessage(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// 发送给发送者和接收者
	sender := h.onlineUsers[message.Username]
	receiver := h.onlineUsers[message.ToUser]

	if sender != nil {
		select {
		case sender.send <- message:
		default:
			close(sender.send)
			delete(h.clients, sender)
		}
	}

	if receiver != nil {
		select {
		case receiver.send <- message:
		default:
			close(receiver.send)
			delete(h.clients, receiver)
		}
	}
}

func (h *Hub) broadcastUserList() {
	h.mutex.RLock()
	users := make([]OnlineUser, 0, len(h.onlineUsers))
	for username, client := range h.onlineUsers {
		users = append(users, OnlineUser{
			Username: username,
			Room:     client.room,
			ID:       client.id,
		})
	}
	h.mutex.RUnlock()

	userListMsg := Message{
		Type:      "userList",
		Content:   string(mustMarshal(users)),
		Timestamp: time.Now(),
		ChatType:  "system",
	}

	h.mutex.RLock()
	for client := range h.clients {
		select {
		case client.send <- userListMsg:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
	h.mutex.RUnlock()
}

func (h *Hub) sendUserListToClient(client *Client) {
	h.mutex.RLock()
	users := make([]OnlineUser, 0, len(h.onlineUsers))
	for username, c := range h.onlineUsers {
		if username != client.username { // 不包括自己
			users = append(users, OnlineUser{
				Username: username,
				Room:     c.room,
				ID:       c.id,
			})
		}
	}
	h.mutex.RUnlock()

	userListMsg := Message{
		Type:      "userList",
		Content:   string(mustMarshal(users)),
		Timestamp: time.Now(),
		ChatType:  "system",
	}

	select {
	case client.send <- userListMsg:
	default:
		close(client.send)
		delete(h.clients, client)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		message.Username = c.username
		message.Timestamp = time.Now()
		message.Source = "websocket"

		// 根据消息类型路由
		if message.ChatType == "private" {
			hub.privateMsg <- message
		} else {
			message.Room = c.room
			message.ChatType = "group"
			message.Type = "message"
			hub.broadcast <- message
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	username := r.URL.Query().Get("username")
	room := r.URL.Query().Get("room")

	if username == "" {
		username = "Anonymous"
	}
	if room == "" {
		room = "general"
	}

	client := &Client{
		conn:     conn,
		send:     make(chan Message, 256),
		username: username,
		room:     room,
		id:       generateClientID(),
	}

	hub.register <- client

	// 发送用户列表给新客户端
	hub.userList <- client

	go client.writePump()
	go client.readPump()
}

func handleUserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	hub.mutex.RLock()
	users := make([]OnlineUser, 0, len(hub.onlineUsers))
	for username, client := range hub.onlineUsers {
		users = append(users, OnlineUser{
			Username: username,
			Room:     client.room,
			ID:       client.id,
		})
	}
	hub.mutex.RUnlock()

	json.NewEncoder(w).Encode(users)
}

func messageToJSON(msg Message) string {
	data, _ := json.Marshal(msg)
	return string(data)
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().UnixNano()%1000))
}

func setupMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://14.103.243.153:1883")
	opts.SetClientID("chat-server")
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		var message Message
		if err := json.Unmarshal(msg.Payload(), &message); err == nil {
			message.Source = "mqtt"

			// 根据主题判断消息类型
			topic := msg.Topic()
			if len(topic) > 11 && topic[:11] == "chat/group/" {
				// 群聊消息
				hub.broadcast <- message
			} else if len(topic) > 13 && topic[:13] == "chat/private/" {
				// 私聊消息
				fmt.Printf(">>>>>>>>>Received private message: %v\n", message)
				hub.privateMsg <- message
			}
		}
	})

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("MQTT connection error: %v", token.Error())
		return
	}

	// 订阅群聊和私聊主题
	if token := mqttClient.Subscribe("chat/group/+", 0, nil); token.Wait() && token.Error() != nil {
		log.Printf("MQTT group subscription error: %v", token.Error())
	}

	if token := mqttClient.Subscribe("chat/private/+/+", 0, nil); token.Wait() && token.Error() != nil {
		log.Printf("MQTT private subscription error: %v", token.Error())
	}

	log.Println("MQTT client connected and subscribed")
}

func main() {
	hub = newHub()
	go hub.run()

	setupMQTT()

	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/ws", handleWebSocket)
	//http.HandleFunc("/api/users", handleUserList)

	log.Println("Chat server starting on :8089")
	log.Fatal(http.ListenAndServe(":8089", nil))
}
