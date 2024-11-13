package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// * RSSIMap defines the structure of the RSSI data
type RSSIMap map[string]int

// * Reads /proc/net/wireless and extracts RSSI for all wireless interfaces.
func getAllRSSI() (RSSIMap, error) {
	file, err := os.Open("/proc/net/wireless")
	if err != nil {
		return nil, fmt.Errorf("could not open /proc/net/wireless: %v", err)
	}
	defer file.Close()

	rssiMap := make(RSSIMap)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNumber++

		// * Skip the first two lines (headers)
		if lineNumber <= 2 {
			continue
		}

		// wlan0: 0000   70.  -40.  -256        0      0      0      0        0        0

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}

			interfaceName := strings.TrimSpace(parts[0])
			fields := strings.Fields(strings.TrimSpace(parts[1]))

			if len(fields) < 4 {
				continue
			}

			// fields[2] typically contains the signal level (RSSI) with a trailing dot, e.g., "-40."
			levelStr := strings.TrimSuffix(fields[2], ".")
			rssi, err := strconv.Atoi(levelStr)
			if err != nil {
				continue
			}

			rssiMap[interfaceName] = rssi
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/net/wireless: %v", err)
	}

	return rssiMap, nil
}

// * Define a WebSocket upgrader with default options
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// * Client represents a connected WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// * Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// * NewHub initializes a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// * Run starts the Hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// * Handles incoming messages from the WebSocket connection
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// * Sends messages from the send channel to the WebSocket connection
func (c *Client) writePump() {
	// * Ping periodically to keep the connection alive
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
				// * The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// * Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// * Send a ping to the client
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// * Handles WebSocket requests from clients
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("Received WebSocket connection request from %s", r.RemoteAddr)

	// * Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	log.Printf("WebSocket connection established with %s", r.RemoteAddr)

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	testData := map[string]string{"message": "connected"}
	data, err := json.Marshal(testData)
	if err != nil {
		log.Printf("Error marshaling test data: %v", err)
	} else {
		client.send <- data
		log.Printf("Sent test message to client: %s", data)
	}

	// * Start goroutines for reading and writing
	go client.writePump()
	go client.readPump(hub)
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()
	addr := ":" + *port

	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client.html")
	})

	// * Start a ticker that fetches RSSI data every second and broadcasts it
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rssiMap, err := getAllRSSI()
				if err != nil {
					log.Printf("Error fetching RSSI: %v", err)
					continue
				}
				// * Serialize the RSSI map to JSON
				data, err := json.Marshal(rssiMap)
				if err != nil {
					log.Printf("Error marshaling RSSI data: %v", err)
					continue
				}

				// * Log the data being broadcasted
				log.Printf("Broadcasting RSSI data: %s", data)

				// * Broadcast the JSON data to all connected clients
				hub.broadcast <- data
			}
		}
	}()

	log.Printf("WebSocket server started on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
