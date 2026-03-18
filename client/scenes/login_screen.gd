extends Control

const SERVER = "http://localhost:8080"

@onready var email_input = $VBoxContainer/EmailInput
@onready var password_input = $VBoxContainer/PasswordInput
@onready var status_label = $VBoxContainer/StatusLabel

func _ready():
	# Centrar la pantalla
	anchor_right = 1.0
	anchor_bottom = 1.0

func _on_login_button_pressed():
	var email = email_input.text.strip_edges()
	var password = password_input.text.strip_edges()

	if email == "" or password == "":
		status_label.text = "❌ Email and password required"
		return

	status_label.text = "Connecting..."
	_login(email, password)

func _on_register_button_pressed():
	var email = email_input.text.strip_edges()
	var password = password_input.text.strip_edges()

	if email == "" or password == "":
		status_label.text = "❌ Email and password required"
		return

	status_label.text = "Registering..."
	_register(email, password)

func _login(email: String, password: String):
	var http = HTTPRequest.new()
	add_child(http)
	http.request_completed.connect(_on_login_response.bind(http))

	var headers = ["Content-Type: application/json"]
	var body = JSON.stringify({"email": email, "password": password})
	http.request(SERVER + "/auth/login", headers, HTTPClient.METHOD_POST, body)

func _register(email: String, password: String):
	print("Registering - email:'", email, "' password:'", password, "'")
	var http = HTTPRequest.new()
	add_child(http)
	http.request_completed.connect(_on_register_response.bind(http))

	var headers = ["Content-Type: application/json"]
	var body = JSON.stringify({"email": email, "password": password})
	http.request(SERVER + "/auth/register", headers, HTTPClient.METHOD_POST, body)

func _on_login_response(_result, response_code, _headers, body, http):
	http.queue_free()
	var text = body.get_string_from_utf8()
	var data = JSON.parse_string(text)

	if data == null:
		status_label.text = "❌ Server error"
		return

	if response_code == 200 and data["success"]:
		if data["message"] == "no_character":
			status_label.text = "✅ Logged in — No character found"
			GameData.account_id = int(data["account_id"])
			get_tree().change_scene_to_file("res://scenes/character_creation.tscn")
			# TODO: pantalla de creación de personaje
		else:
			status_label.text = "✅ Welcome " + data["player_name"] + "!"
			GameData.player_name = data["player_name"]
			GameData.character_id = int(data["character_id"])
			GameData.account_id = int(data["account_id"])
			GameData.position_x = float(data.get("position_x", 0.0))
			GameData.position_y = float(data.get("position_y", 0.0))
			GameData.position_z = float(data.get("position_z", 0.0))
			get_tree().change_scene_to_file("res://scenes/world.tscn")
	else:
		status_label.text = "❌ " + data.get("message", "Login failed")

func _on_register_response(_result, response_code, _headers, body, http):
	http.queue_free()
	var text = body.get_string_from_utf8()
	var data = JSON.parse_string(text)

	if data == null:
		status_label.text = "❌ Server error"
		return

	if response_code == 201 and data["success"]:
		status_label.text = "✅ Account created! Please login."
	else:
		status_label.text = "❌ " + data.get("message", "Register failed")
