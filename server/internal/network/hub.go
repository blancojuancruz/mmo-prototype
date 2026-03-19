package network

import (
	"encoding/json"
	"fmt"
	"log"
	"mmorpg-server/internal/db"
	"mmorpg-server/internal/game"
	"strconv"
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

type AttackMessage struct {
	Type   string `json:"type"`
	NpcID  string `json:"npc_id"`
	Damage int    `json:"damage"`
}

type ConnectingNpcs struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
	MaxLife     int     `json:"max_life"`
	CurrentLife int     `json:"current_life"`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			for _, npc := range h.gameState.Npcs {
				msg, err := json.Marshal(ConnectingNpcs{
					Type:        "npc_spawn",
					ID:          npc.ID,
					X:           npc.ActualPositionX,
					Y:           npc.ActualPositionY,
					Z:           npc.ActualPositionZ,
					MaxLife:     npc.MaxLife,
					CurrentLife: npc.CurrentLife,
				})
				if err != nil {
					log.Println("❌ Error marshaling NPC:", err)
					continue
				}
				client.send <- msg
			}

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
			var msg struct {
				Type   string  `json:"type"`
				ID     string  `json:"id"`
				NpcID  string  `json:"npc_id"`
				Damage int     `json:"damage"`
				X      float64 `json:"x"`
				Y      float64 `json:"y"`
				Z      float64 `json:"z"`
				RotY   float64 `json:"rot_y"`
			}
			if err := json.Unmarshal(message, &msg); err == nil {
				switch msg.Type {
				case "move":
					h.gameState.UpdatePlayer(&game.PlayerState{
						ID:   msg.ID,
						Type: msg.Type,
						X:    msg.X,
						Y:    msg.Y,
						Z:    msg.Z,
						RotY: msg.RotY,
					})
				case "attack":
					if npc, ok := h.gameState.Npcs[msg.NpcID]; ok {
						npc.AttackChan <- game.AttackEvent{
							PlayerId: msg.ID,
							Damage:   msg.Damage,
						}
					}
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

func (h *Hub) LoadNpcs(npcs []db.NpcSpawnFull) {
	for _, spawn := range npcs {
		npc := &game.NpcState{
			ID:              strconv.Itoa(spawn.ID),
			State:           spawn.State,
			MaxLife:         spawn.MaxLife,
			Damage:          spawn.Damage,
			AttackSpeed:     spawn.AttackSpeed,
			SpawnTimer:      spawn.SpawnTimer,
			CurrentLife:     spawn.CurrentLife,
			ActualPositionX: spawn.ActualPositionX,
			ActualPositionY: spawn.ActualPositionY,
			ActualPositionZ: spawn.ActualPositionZ,
			SpawnPositionX:  spawn.SpawnPositionX,
			SpawnPositionY:  spawn.SpawnPositionY,
			SpawnPositionZ:  spawn.SpawnPositionZ,
		}
		npc.AttackChan = make(chan game.AttackEvent, 10)
		npc.BroadcastChan = h.broadcast
		h.gameState.AddNpc(npc)
	}
}
