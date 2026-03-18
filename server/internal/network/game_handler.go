package network

import (
	"encoding/json"
	"net/http"

	"mmorpg-server/internal/db"

	"github.com/jmoiron/sqlx"
)

type SavePositionRequest struct {
	CharacterID int     `json:"character_id"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
}

func SavePositionHandler(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SavePositionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, map[string]bool{"success": false}, http.StatusBadRequest)
			return
		}

		err := db.SaveCharacterPosition(database, req.CharacterID, req.X, req.Y, req.Z)
		if err != nil {
			sendJSON(w, map[string]bool{"success": false}, http.StatusInternalServerError)
			return
		}

		sendJSON(w, map[string]bool{"success": true}, http.StatusOK)
	}
}
