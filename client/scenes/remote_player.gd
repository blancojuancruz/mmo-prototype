extends Node3D

@onready var anim_player = $Character/AnimationPlayer

var target_position = Vector3.ZERO
var target_rotation_y = 0.0

func _ready():
	anim_player.play("AnimPack1/Idle")

func update_state(data: Dictionary):
	target_position = Vector3(data["x"], data["y"], data["z"])
	target_rotation_y = data["rot_y"]

func _physics_process(delta):
	# Interpolación suave para movimiento fluido
	position = position.lerp(target_position, delta * 10)
	rotation.y = lerp_angle(rotation.y, target_rotation_y, delta * 10)

	# Animación basada en velocidad
	var speed = (target_position - position).length()
	if speed > 0.01:
		if anim_player.current_animation != "AnimPack1/Walk":
			anim_player.play("AnimPack1/Walk")
	else:
		if anim_player.current_animation != "AnimPack1/Idle":
			anim_player.play("AnimPack1/Idle")
