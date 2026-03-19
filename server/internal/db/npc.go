package db

import "github.com/jmoiron/sqlx"

type Npc struct {
	Name        string `db:"name"`
	DamageType  string `db:"damage_type"` // Future FK Table Npc_Damage_Types
	RangeType   string `db:"range_type"`  // Future FK Table Npc_Range_Type
	ID          int    `db:"id"`          // Unique
	MaxLife     int    `db:"max_life"`
	Level       int    `db:"level"`
	Damage      int    `db:"damage"`
	AttackSpeed int    `db:"attack_speed"`
}

func CreateNpc(db *sqlx.DB, name, damageType, rangeType string, maxLife, level, damage, attackSpeed int) (*Npc, error) {
	npc := &Npc{}
	err := db.QueryRowx(
		`INSERT INTO npcs (name, max_life, level, damage, damage_type, range_type, attack_speed)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, name, max_life, level, damage, damage_type, range_type, attack_speed
		`, name, maxLife, level, damage, damageType, rangeType, attackSpeed).StructScan(npc)
	return npc, err
}
