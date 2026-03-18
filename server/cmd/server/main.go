package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"mmorpg-server/internal/db"
	"mmorpg-server/internal/network"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🎮 MMORPG Server starting...")

	// Cargar .env — si no existe usa las variables del sistema
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("❌ Invalid DB_PORT value:", err)
	}
	database := db.NewPostgres(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})

	db.RunMigrations(database)

	hub := network.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		network.ServeWS(hub, w, r)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("pong")); err != nil {
			log.Println("❌ Error writing response:", err)
		}
	})
	http.HandleFunc("/auth/register", network.RegisterHandler(database))
	http.HandleFunc("/auth/login", network.LoginHandler(database))
	http.HandleFunc("/auth/character", network.CreateCharacterHandler(database))
	http.HandleFunc("/game/save_position", network.SavePositionHandler(database))

	serverPort := os.Getenv("SERVER_PORT")
	fmt.Printf("Server listening on port :%s\n", serverPort)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		database.Close()
		log.Fatal("❌ Server error:", err)
	}
}
