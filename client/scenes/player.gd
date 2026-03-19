extends CharacterBody3D

const SPEED = 5.0
const GRAVITY = 9.8
const MOUSE_SENSITIVITY = 0.005
const CAMERA_MIN_ANGLE = -60.0
const CAMERA_MAX_ANGLE = 20.0
const SERVER_URL = "ws://localhost:8080/ws"
const SEND_RATE = 0.05
const SAVE_RATE = 30.0

@onready var spring_arm = $SpringArm3D
@onready var anim_player = $Character/AnimationPlayer

var camera_rotation_x = -20.0
var socket = WebSocketPeer.new()
var player_id = GameData.player_name
var connected = false
var send_timer = 0.0
var remote_players = {}
var remote_player_scene = preload("res://scenes/remote_player.tscn")
var is_quitting = false
var save_timer = 0.0

func _ready():
	Input.set_mouse_mode(Input.MOUSE_MODE_CAPTURED)
	spring_arm.rotation_degrees.x = camera_rotation_x
	anim_player.play("AnimPack1/Idle")
	_connect_to_server()
	get_tree().set_auto_accept_quit(false)
	position = Vector3(GameData.position_x, GameData.position_y, GameData.position_z)

func _connect_to_server():
	var err = socket.connect_to_url(SERVER_URL + "?id=" + player_id)
	if err != OK:
		print("❌ Error connecting to server")
	else:
		print("🔌 Connecting to server...")

func _input(event):
	if event is InputEventMouseMotion:
		rotation.y -= event.relative.x * MOUSE_SENSITIVITY
		camera_rotation_x -= event.relative.y * MOUSE_SENSITIVITY * 57.3
		camera_rotation_x = clamp(camera_rotation_x, CAMERA_MIN_ANGLE, CAMERA_MAX_ANGLE)
		spring_arm.rotation_degrees.x = camera_rotation_x
	if event.is_action_pressed("ui_cancel"):
		Input.set_mouse_mode(Input.MOUSE_MODE_VISIBLE)

func _physics_process(delta):
	save_timer += delta
	if save_timer >= SAVE_RATE:
		save_timer = 0.0
		_save_position()
	socket.poll()
	var state = socket.get_ready_state()

	if state == WebSocketPeer.STATE_OPEN:
		if not connected:
			connected = true
			print("✅ Connected as: " + player_id)
		send_timer += delta
		if send_timer >= SEND_RATE:
			send_timer = 0.0
			_send_position()
		_receive_messages()
	elif state == WebSocketPeer.STATE_CLOSED and connected:
		connected = false
		print("❌ Disconnected")

	if not is_on_floor():
		velocity.y -= GRAVITY * delta

	var direction = Vector3.ZERO
	if Input.is_action_pressed("move_forward") or Input.is_action_pressed("ui_up"):
		direction -= transform.basis.z
	if Input.is_action_pressed("move_back") or Input.is_action_pressed("ui_down"):
		direction += transform.basis.z
	if Input.is_action_pressed("move_left") or Input.is_action_pressed("ui_left"):
		direction -= transform.basis.x
	if Input.is_action_pressed("move_right") or Input.is_action_pressed("ui_right"):
		direction += transform.basis.x

	var is_moving = direction != Vector3.ZERO
	var is_running = Input.is_action_pressed("run")

	if is_moving:
		direction = direction.normalized()
		var current_speed = SPEED * 2.0 if is_running else SPEED
		velocity.x = direction.x * current_speed
		velocity.z = direction.z * current_speed
		if is_running:
			if anim_player.current_animation != "AnimPack1/Run":
				anim_player.play("AnimPack1/Run")
		else:
			if anim_player.current_animation != "AnimPack1/Walk":
				anim_player.play("AnimPack1/Walk")
	else:
		velocity.x = move_toward(velocity.x, 0, SPEED)
		velocity.z = move_toward(velocity.z, 0, SPEED)
		if anim_player.current_animation != "AnimPack1/Idle":
			anim_player.play("AnimPack1/Idle")

	move_and_slide()

func _send_position():
	var data = {
		"id": player_id,
		"type": "move",
		"x": position.x,
		"y": position.y,
		"z": position.z,
		"rot_y": rotation.y
	}
	socket.send_text(JSON.stringify(data))

func _receive_messages():
	var world = get_parent()
	while socket.get_available_packet_count() > 0:
		var packet = socket.get_packet()
		var text = packet.get_string_from_utf8()
		var data = JSON.parse_string(text)
		if data == null:
			return
		var msg_id = data.get("id", "")
		var msg_type = data.get("type", "")
		if msg_type == "npc_spawn":
			world.spawn_npc(data)
		if msg_id == player_id:
			return
		if msg_type == "move":
			_handle_remote_move(data)
		elif msg_type == "disconnect":
			_handle_remote_disconnect(msg_id)
		elif msg_type == "chat":
			world.add_chat_message(msg_id, data.get("message", ""))

func _handle_remote_move(data: Dictionary):
	var id = data["id"]
	if not remote_players.has(id):
		# Instanciar nuevo jugador remoto
		var remote = remote_player_scene.instantiate()
		remote.name = id
		get_parent().add_child(remote)
		remote_players[id] = remote
		print("👤 New player joined: " + id)
	remote_players[id].update_state(data)

func _handle_remote_disconnect(id: String):
	if remote_players.has(id):
		remote_players[id].queue_free()
		remote_players.erase(id)
		print("👋 Player left: " + id)

func _notification(what): 
	if what == NOTIFICATION_WM_CLOSE_REQUEST:
		_save_position()
		await get_tree().create_timer(0.5).timeout
		get_tree().quit()

func _save_position(): 
	var http = HTTPRequest.new()
	add_child(http)
	http.request_completed.connect(_on_position_saved.bind(http).bind(http))
	var headers = ["Content-Type: application/json"]
	var body = JSON.stringify({
		"character_id": GameData.character_id,
		"x": position.x,
		"y": position.y,
		"z": position.z,
	})
	http.request("http://localhost:8080/game/save_position", headers, HTTPClient.METHOD_POST, body)

func _on_position_saved(_result, _code, _headers, _body, http):
	http.queue_free()
	get_tree().quit()

func send_chat(message: String):
	var data = {
		"id": player_id,
		"type": "chat",
		"message": message
	}
	socket.send_text(JSON.stringify(data))
	var world = get_parent()
	world.add_chat_message(player_id, message)
