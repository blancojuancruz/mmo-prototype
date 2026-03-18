package network

import (
	"encoding/json"
	"fmt"
	"log"
	"mmorpg-server/internal/game"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	gameState  *game.GameState
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		gameState:  game.NewGameState(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("✅ Player connected: %s. Total: %d\n", client.ID, len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.gameState.RemovePlayer(client.ID)

				disconnectMsg, err := json.Marshal(map[string]string{
					"type": "disconnect",
					"id":   client.ID,
				})
				if err != nil {
					log.Println("❌ Error marshaling disconnect message:", err)
					return
				}
				h.broadcastToAll(disconnectMsg)
				fmt.Printf("❌ Player disconnected: %s. Total: %d\n", client.ID, len(h.clients))
			}

		case message := <-h.broadcast:
			var state game.PlayerState
			if err := json.Unmarshal(message, &state); err == nil {
				if state.Type == "move" {
					h.gameState.UpdatePlayer(&state)
				}
			}
			h.broadcastToAll(message)
		}
	}
}

func (h *Hub) broadcastToAll(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
