package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			id         SERIAL PRIMARY KEY,
			email      VARCHAR(255) UNIQUE NOT NULL,
			password   VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS character_classes (
			id          SERIAL PRIMARY KEY,
			name        VARCHAR(50) UNIQUE NOT NULL,
			description TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS zones (
			id           SERIAL PRIMARY KEY,
			name         VARCHAR(50) UNIQUE NOT NULL,
			display_name VARCHAR(100) NOT NULL,
			max_players  INT DEFAULT 100
		)`,

		`CREATE TABLE IF NOT EXISTS characters (
			id          SERIAL PRIMARY KEY,
			account_id  INT NOT NULL REFERENCES accounts(id),
			name        VARCHAR(50) UNIQUE NOT NULL,
			class_id    INT NOT NULL REFERENCES character_classes(id),
			level       INT DEFAULT 1,
			experience  BIGINT DEFAULT 0,
			position_x  FLOAT DEFAULT 0,
			position_y  FLOAT DEFAULT 0,
			position_z  FLOAT DEFAULT 0,
			zone_id     INT REFERENCES zones(id),
			created_at  TIMESTAMP DEFAULT NOW()
		)`,

		`INSERT INTO character_classes (name, description)
		VALUES
			('Warrior', 'Melee fighter with high defense'),
			('Mage', 'Spellcaster with high damage'),
			('Archer', 'Ranged fighter with high speed')
		ON CONFLICT DO NOTHING`,

		`INSERT INTO zones (name, display_name, max_players)
		VALUES
			('starting_zone', 'Plains of Beginning', 100),
			('forest_zone', 'Ancient Forest', 50),
			('dungeon_01', 'Cave of Trials', 20)
		ON CONFLICT DO NOTHING`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			log.Fatal("❌ Migration failed:", err)
		}
	}

	log.Println("✅ Migrations completed")
}