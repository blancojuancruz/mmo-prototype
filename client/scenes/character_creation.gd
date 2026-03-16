extends Control

const SERVER = "http://localhost:8080"

@onready var name_input = $VBoxContainer/NameInput
@onready var class_option = $VBoxContainer/ClassOption
@onready var status_label = $VBoxContainer/StatusLabel

func _ready():
	anchor_right = 1.0
	anchor_bottom = 1.0
	# Agregar las clases disponibles
	class_option.add_item("Warrior", 1)
	class_option.add_item("Mage", 2)
	class_option.add_item("Archer", 3)

func _on_create_button_pressed():
	var char_name = name_input.text.strip_edges()
	if char_name == "":
		status_label.text = "❌ Name required"
		return

	status_label.text = "Creating character..."

	var http = HTTPRequest.new()
	add_child(http)
	http.request_completed.connect(_on_create_response.bind(http))

	var headers = ["Content-Type: application/json"]
	# get_selected_id() devuelve el ID que asignamos en add_item()
	var class_id = class_option.get_selected_id()
	var body = JSON.stringify({
		"account_id": GameData.account_id,
		"name": char_name,
		"class_id": class_id
	})
	http.request(SERVER + "/auth/character", headers, HTTPClient.METHOD_POST, body)

func _on_create_response(_result, response_code, _headers, body, http):
	http.queue_free()
	var text = body.get_string_from_utf8()
	var data = JSON.parse_string(text)

	if data == null:
		status_label.text = "❌ Server error"
		return

	if response_code == 201 and data["success"]:
		GameData.player_name = data["player_name"]
		GameData.character_id = int(data["character_id"])
		get_tree().change_scene_to_file("res://scenes/world.tscn")
	else:
		status_label.text = "❌ " + data.get("message", "Failed to create character")
