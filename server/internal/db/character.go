package db

import "github.com/jmoiron/sqlx"

type Character struct {
	ID         int     `db:"id"`
	AccountID  int     `db:"account_id"`
	Name       string  `db:"name"`
	ClassID    int     `db:"class_id"`
	Level      int     `db:"level"`
	Experience int64   `db:"experience"`
	PositionX  float64 `db:"position_x"`
	PositionY  float64 `db:"position_y"`
	PositionZ  float64 `db:"position_z"`
	ZoneID     int     `db:"zone_id"`
}

func CreateCharacter(db *sqlx.DB, accountID int, name string, classID int) (*Character, error) {
	character := &Character{}
	err := db.QueryRowx(
		`INSERT INTO characters (account_id, name, class_id, zone_id)
		 VALUES ($1, $2, $3, 1)
		 RETURNING id, account_id, name, class_id, level, experience,
		           position_x, position_y, position_z, zone_id`,
		accountID, name, classID,
	).StructScan(character)
	return character, err
}

func GetCharacterByAccountID(db *sqlx.DB, accountID int) (*Character, error) {
	character := &Character{}
	err := db.QueryRowx(
		`SELECT id, account_id, name, class_id, level, experience,
		        position_x, position_y, position_z, zone_id
		 FROM characters WHERE account_id = $1 LIMIT 1`,
		accountID,
	).StructScan(character)
	return character, err
}

func SaveCharacterPosition(db *sqlx.DB, characterID int, x, y, z float64) error {
	_, err := db.Exec(
		`UPDATE characters SET position_x=$1, position_y=$2, position_z=$3
		 WHERE id=$4`,
		x, y, z, characterID,
	)
	return err
}