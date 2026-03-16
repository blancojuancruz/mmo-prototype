package network

import (
	"encoding/json"
	"net/http"

	"mmorpg-server/internal/db"

	"github.com/jmoiron/sqlx"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
	AccountID   int     `json:"account_id,omitempty"`
	CharacterID int     `json:"character_id,omitempty"`
	PlayerName  string  `json:"player_name,omitempty"`
	PositionX   float64 `json:"position_x,omitempty"`
	PositionY   float64 `json:"position_y,omitempty"`
	PositionZ   float64 `json:"position_z,omitempty"`
}

func RegisterHandler(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, AuthResponse{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			sendJSON(w, AuthResponse{Success: false, Message: "Email and password required"}, http.StatusBadRequest)
			return
		}

		account, err := db.CreateAccount(database, req.Email, req.Password)
		if err != nil {
			sendJSON(w, AuthResponse{Success: false, Message: "Email already exists"}, http.StatusConflict)
			return
		}

		sendJSON(w, AuthResponse{
			Success:   true,
			Message:   "Account created successfully",
			AccountID: account.ID,
		}, http.StatusCreated)
	}
}

func LoginHandler(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, AuthResponse{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}

		account, err := db.GetAccountByEmail(database, req.Email)
		if err != nil {
			sendJSON(w, AuthResponse{Success: false, Message: "Invalid credentials"}, http.StatusUnauthorized)
			return
		}

		if !db.ValidatePassword(account, req.Password) {
			sendJSON(w, AuthResponse{Success: false, Message: "Invalid credentials"}, http.StatusUnauthorized)
			return
		}

		character, err := db.GetCharacterByAccountID(database, account.ID)
		if err != nil {
			sendJSON(w, AuthResponse{
				Success:   true,
				Message:   "no_character",
				AccountID: account.ID,
			}, http.StatusOK)
			return
		}

		sendJSON(w, AuthResponse{
      Success:     true,
      Message:     "ok",
      AccountID:   account.ID,
      CharacterID: character.ID,
      PlayerName:  character.Name,
      PositionX:   character.PositionX,
      PositionY:   character.PositionY,
      PositionZ:   character.PositionZ,
    }, http.StatusOK)
	}
}

func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type CreateCharacterRequest struct {
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	ClassID   int    `json:"class_id"`
}

func CreateCharacterHandler(database *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateCharacterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, AuthResponse{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.AccountID == 0 || req.ClassID == 0 {
			sendJSON(w, AuthResponse{Success: false, Message: "Name, account_id and class_id required"}, http.StatusBadRequest)
			return
		}

		character, err := db.CreateCharacter(database, req.AccountID, req.Name, req.ClassID)
		if err != nil {
			sendJSON(w, AuthResponse{Success: false, Message: "Name already taken"}, http.StatusConflict)
			return
		}

		sendJSON(w, AuthResponse{
			Success:     true,
			Message:     "ok",
			AccountID:   req.AccountID,
			CharacterID: character.ID,
			PlayerName:  character.Name,
		}, http.StatusCreated)
	}
}