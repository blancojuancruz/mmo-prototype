package db

import "github.com/jmoiron/sqlx"

type NpcSpawn struct {
	State           string  `db:"state"` // idle/combat/dead
	ActualPositionX float64 `db:"actual_position_x"`
	ActualPositionY float64 `db:"actual_position_y"`
	ActualPositionZ float64 `db:"actual_position_z"`
	SpawnPositionX  float64 `db:"spawn_position_x"`
	SpawnPositionY  float64 `db:"spawn_position_y"`
	SpawnPositionZ  float64 `db:"spawn_position_z"`
	CurrentLife     int     `db:"state"`
	ID              int     `db:"id"`
	NpcId           int     `db:"npc_id"`  // FK Table Npcs
	ZoneId          int     `db:"zone_id"` // FK Table Zones
	SpawnTimer      int     `db:"spawn_timer"`
}

type NpcSpawnFull struct {
	State           string  `db:"state"`
	ActualPositionX float64 `db:"actual_position_x"`
	ActualPositionY float64 `db:"actual_position_y"`
	ActualPositionZ float64 `db:"actual_position_z"`
	SpawnPositionX  float64 `db:"spawn_position_x"`
	SpawnPositionY  float64 `db:"spawn_position_y"`
	SpawnPositionZ  float64 `db:"spawn_position_z"`
	AttackSpeed     float32 `db:"attack_speed"`
	ID              int     `db:"id"`
	CurrentLife     int     `db:"current_life"`
	MaxLife         int     `db:"max_life"`
	Damage          int     `db:"damage"`
	SpawnTimer      int     `db:"spawn_timer"`
}

func GetAllNpcSpawns(db *sqlx.DB) ([]NpcSpawnFull, error) {
	spawns := []NpcSpawnFull{}
	err := db.Select(
		&spawns,
		`
			SELECT npcs_spawns.id, npcs.max_life, npcs.damage, npcs.attack_speed, npcs_spawns.state, npcs_spawns.current_life, npcs_spawns.spawn_timer, 
			npcs_spawns.actual_position_x, npcs_spawns.actual_position_y, npcs_spawns.actual_position_z, npcs_spawns.spawn_position_x, npcs_spawns.spawn_position_y, npcs_spawns.spawn_position_z
			FROM npcs
			INNER JOIN npcs_spawns ON npcs.id = npcs_spawns.npc_id`,
	)

	return spawns, err
}

func CreateNpcSpawn(db *sqlx.DB, state string, npcId, zoneId, currentLife, spawnTimer int, actualPositionX, actualPositionY, actualPositionZ, spawnPositionX, spawnPositionY, spawnPositionZ float64) (*NpcSpawn, error) {
	npcSpawn := &NpcSpawn{}
	err := db.QueryRowx(
		`INSERT INTO npcs_spawns (state, npc_id, zone_id, current_life, spawn_timer, actual_position_X, actual_position_Y, actual_position_Z, spawn_position_X, spawn_position_Y, spawn_position_Z)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, state, npc_id, zone_id, current_life, spawn_timer, actual_position_X, actual_position_Y, actual_position_Z, spawn_position_X, spawn_position_Y, spawn_position_Z
		`, state, npcId, zoneId, currentLife, spawnTimer, actualPositionX, actualPositionY, actualPositionZ, spawnPositionX, spawnPositionY, spawnPositionZ).StructScan(npcSpawn)
	return npcSpawn, err
}
