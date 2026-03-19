package game

import (
	"sync"
)

type GameState struct {
	Players map[string]*PlayerState
	Npcs    map[string]*NpcState
	mu      sync.RWMutex
}

func NewGameState() *GameState {
	return &GameState{
		Players: make(map[string]*PlayerState),
		Npcs:    make(map[string]*NpcState),
	}
}

func (gs *GameState) UpdatePlayer(state *PlayerState) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Players[state.ID] = state
}

func (gs *GameState) RemovePlayer(id string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	delete(gs.Players, id)
}

func (gs *GameState) AddNpc(npc *NpcState) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Npcs[npc.ID] = npc
	go npc.Run()
}

func (gs *GameState) UpdateNpc(state *NpcState) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Npcs[state.ID] = state
}

func (gs *GameState) RemoveNpc(id string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	delete(gs.Npcs, id)
}
