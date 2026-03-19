extends StaticBody3D

@onready var health_label = $HealthLabel

var npc_id: int = 0
var npc_state = "idle"
var current_life: int = 0
var max_life: int = 0

func update_health(new_life: int):
	current_life = new_life
	health_label.text = str(current_life) + "/" + str(max_life)
	
func take_damage(damage):
	var new_life = current_life - damage
	update_health(new_life)
	if current_life <= 0:
		npc_state = "dead"
