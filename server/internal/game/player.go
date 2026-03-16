package game

import "sync"

type PlayerState struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	RotY  float64 `json:"rot_y"`
}

type GameState struct {
	mu      sync.RWMutex
	Players map[string]*PlayerState
}

func NewGameState() *GameState {
	return &GameState{
		Players: make(map[string]*PlayerState),
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