extends Node3D

@onready var chat_messages = $ChatUI/VBoxContainer/ScrollContainer/ChatMessages
@onready var chat_input = $ChatUI/VBoxContainer/HBoxContainer/ChatInput

var remote_npcs = {}
var remote_npc_scene = preload("res://scenes/npc.tscn")

func _ready():
	chat_input.gui_input.connect(_on_chat_input)

func _on_chat_input(event):
	if event is InputEventKey and event.pressed:
		if event.keycode == KEY_ENTER and chat_input.text.strip_edges() != "":
			_send_chat()

func _on_send_button_pressed():
	if chat_input.text.strip_edges() != "":
		_send_chat()

func _send_chat():
	var message = chat_input.text.strip_edges()
	chat_input.text = ""
	# Acceder al Player y enviar via WebSocket
	var player = $Player
	player.send_chat(message)

func add_chat_message(sender: String, message: String):
	var label = Label.new()
	label.text = sender + ": " + message
	label.autowrap_mode = TextServer.AUTOWRAP_WORD
	chat_messages.add_child(label)

	# Auto-scroll hacia abajo
	await get_tree().process_frame
	var scroll = $ChatUI/VBoxContainer/ScrollContainer
	scroll.scroll_vertical = scroll.get_v_scroll_bar().max_value

func spawn_npc(data):
	var remote = remote_npc_scene.instantiate()
	remote.position = Vector3(data["x"], 1.0, data["z"])
	remote.npc_spawn_id = str(data["id"])
	remote.max_life = data["max_life"]
	remote.current_life = data["current_life"]
	remote_npcs[data["id"]] = remote
	add_child(remote)

func _unhandled_input(event):
	if event is InputEventMouseButton:
		if event.button_index == MOUSE_BUTTON_RIGHT:
			if event.is_pressed():
				var max_distance = 150
				var click_origin
				var click_direction
				var click_destiny
				var world_physics = get_world_3d().direct_space_state
				var camera = get_viewport().get_camera_3d()
				click_origin = camera.project_ray_origin(event.position)
				click_direction = camera.project_ray_normal(event.position)
				click_destiny = click_origin + click_direction * max_distance
				var query = PhysicsRayQueryParameters3D.create(click_origin, click_destiny)
				var result = world_physics.intersect_ray(query)
				if !result.is_empty():
					if result["collider"].is_in_group("npc"):
						var player = $Player
						player.attack_npc(result["collider"].npc_spawn_id)
