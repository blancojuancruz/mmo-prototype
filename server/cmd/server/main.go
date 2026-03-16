package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"mmorpg-server/internal/db"
	"mmorpg-server/internal/network"
)

func main() {
	fmt.Println("🎮 MMORPG Server starting...")

	// Cargar .env — si no existe usa las variables del sistema
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	database := db.NewPostgres(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	})
	defer database.Close()

	db.RunMigrations(database)

	hub := network.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		network.ServeWS(hub, w, r)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	http.HandleFunc("/auth/register", network.RegisterHandler(database))
	http.HandleFunc("/auth/login", network.LoginHandler(database))
	http.HandleFunc("/auth/character", network.CreateCharacterHandler(database))
	http.HandleFunc("/game/save_position", network.SavePositionHandler(database))

	serverPort := os.Getenv("SERVER_PORT")
	fmt.Printf("Server listening on port :%s\n", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}