package game

import (
	"encoding/json"
	"log"
	"time"
)

type NpcState struct {
	AttackChan      chan AttackEvent
	BroadcastChan   chan []byte
	ID              string
	State           string
	TargetPlayerId  string
	ActualPositionX float64
	ActualPositionY float64
	ActualPositionZ float64
	SpawnPositionX  float64
	SpawnPositionY  float64
	SpawnPositionZ  float64
	AttackTimer     float64
	AttackSpeed     float32
	CurrentLife     int
	MaxLife         int
	Damage          int
	SpawnTimer      int
}

type AttackEvent struct {
	PlayerId string
	Damage   int
}

func (n *NpcState) Run() {
	for {
		switch n.State {
		case "idle":
			event := <-n.AttackChan
			n.State = "combat"
			n.TargetPlayerId = event.PlayerId
			n.CurrentLife -= event.Damage
		case "combat":
			select {
			case <-time.After(time.Duration(n.AttackSpeed) * time.Second):
				n.broadcastDamage()
			case event := <-n.AttackChan:
				n.CurrentLife -= event.Damage

				if n.CurrentLife <= 0 {
					n.State = "dead"
				}

				if n.CurrentLife > 0 && n.AttackTimer <= 0 {
					n.broadcastDamage()
				}

			}
		case "dead":
			<-time.After(time.Duration(n.SpawnTimer) * time.Second)
			n.CurrentLife = n.MaxLife
			n.ActualPositionX = n.SpawnPositionX
			n.ActualPositionY = n.SpawnPositionY
			n.ActualPositionZ = n.SpawnPositionZ
			n.State = "idle"
		}
	}
}

func (n *NpcState) broadcastDamage() {
	combatMsg, err := json.Marshal(map[string]any{
		"type":   "player_damage",
		"target": n.TargetPlayerId,
		"damage": n.Damage,
	})

	if err != nil {
		log.Println("❌ Error generating combat log:", err)
		return
	}
	n.BroadcastChan <- combatMsg
}
