extends Node3D

@onready var chat_messages = $ChatUI/VBoxContainer/ScrollContainer/ChatMessages
@onready var chat_input = $ChatUI/VBoxContainer/HBoxContainer/ChatInput

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
