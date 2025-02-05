package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Structure to hold messages for each room
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Image    string `json:"image"`
}

var (
	rooms    = make(map[string]map[*websocket.Conn]bool)
	messages = make(map[string][]Message)
	clients  = make(map[*websocket.Conn]bool) // Track all connected clients for broadcasting new rooms
	mu       sync.Mutex
)

func main() {
	setupRoutes()
	startServer()
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws/", handleConnections)
	http.HandleFunc("/create-room", handleRoomCreation) // New route for handling room creation
}

func startServer() {
	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	roomName := getRoomName(r)
	if roomName == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	conn, err := upgradeConnection(w, r)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer closeConnection(roomName, conn)

	registerConnection(roomName, conn)
	sendExistingMessages(roomName, conn)

	listenForMessages(roomName, conn)
}

func getRoomName(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/ws/")
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		mu.Lock()
		clients[conn] = true // Add to the global client list
		mu.Unlock()
	}
	return conn, err
}

func closeConnection(roomName string, conn *websocket.Conn) {
	conn.Close()
	mu.Lock()
	delete(rooms[roomName], conn)
	delete(clients, conn) // Remove from the global client list
	mu.Unlock()
	fmt.Printf("User disconnected from room: %s\n", roomName)
}

func registerConnection(roomName string, conn *websocket.Conn) {
	mu.Lock()
	if rooms[roomName] == nil {
		rooms[roomName] = make(map[*websocket.Conn]bool)
	}
	rooms[roomName][conn] = true
	mu.Unlock()
}

func sendExistingMessages(roomName string, conn *websocket.Conn) {
	mu.Lock()
	fmt.Printf("User joined room: %s. Sending existing messages...\n", roomName)
	for _, msg := range messages[roomName] {
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("Error sending message history:", err)
		}
	}
	mu.Unlock()
}

func listenForMessages(roomName string, conn *websocket.Conn) {
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			handleReadError(err, roomName)
			break
		}
		storeAndBroadcastMessage(roomName, msg)
	}
}

func handleReadError(err error, roomName string) {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		fmt.Println("Unexpected WebSocket closure:", err)
	} else {
		fmt.Println("User disconnected from room:", roomName)
	}
}

func storeAndBroadcastMessage(roomName string, msg Message) {
	mu.Lock()
	messages[roomName] = append(messages[roomName], msg)
	fmt.Printf("Stored message in room: %s. Total messages: %d\n", roomName, len(messages[roomName]))
	broadcastMessage(roomName, msg)
	mu.Unlock()
}

func broadcastMessage(roomName string, msg Message) {
	for conn := range rooms[roomName] {
		if err := conn.WriteJSON(msg); err != nil {
			fmt.Println("Error broadcasting message:", err)
			conn.Close()
			delete(rooms[roomName], conn)
		}
	}
}

// Handle room creation and broadcast to all clients
func handleRoomCreation(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection for room creation:", err)
		return
	}
	defer conn.Close()

	for {
		var roomName string
		err := conn.ReadJSON(&roomName)
		if err != nil {
			fmt.Println("Error reading room name:", err)
			break
		}

		mu.Lock()
		if rooms[roomName] == nil {
			rooms[roomName] = make(map[*websocket.Conn]bool)
			broadcastNewRoom(roomName)
		}
		mu.Unlock()
	}
}

func broadcastNewRoom(roomName string) {
	roomUpdate := map[string]string{"newRoom": roomName}
	for conn := range clients {
		if err := conn.WriteJSON(roomUpdate); err != nil {
			fmt.Println("Error broadcasting new room:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
