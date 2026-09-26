package gb

type Joypad struct {
	sel     byte // bits 4 and 5 in P1 (select d-pad and select buttons)
	dpad    byte // bits 0, 1, 2, 3 for right, left, up, down
	buttons byte // bits 0, 1, 2, 3 for A, B, start, select
}

func InitJoypad() *Joypad {
	return &Joypad{
		dpad:    0x0,
		buttons: 0x0,
		sel:     0x30,
	}
}

func (j *Joypad) Read() byte {
	var low byte = 0xf
	if !bitEnabled(j.sel, 4) {
		low &= ^j.dpad
	}

	if !bitEnabled(j.sel, 5) {
		low &= ^j.buttons
	}

	return 0xc0 | j.sel | low
}

func (j *Joypad) Write(val byte) {
	// only bits 4 and 5, which we represent in sel, are writable
	j.sel = val & 0x30
}

// Update determines whether a button has been pressed this frame, determining whether
// a joypad interrupt needs to be requested
func (j *Joypad) Update(inp byte) bool {
	// if we derive that anything has been pressed, we return a JoypadInterruptRequest
	oldJoypad := j.Read()
	j.buttons = (inp & 0x0f)
	j.dpad = (inp >> 4)
	newJoypad := j.Read()

	return oldJoypad&^newJoypad&0x0f != 0
}
