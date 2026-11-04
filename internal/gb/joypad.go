package gb

type Joypad struct {
	sel     byte // bits 4 and 5 in P1 (select d-pad and select buttons)
	dpad    byte // bits 0, 1, 2, 3 for right, left, up, down
	buttons byte // bits 0, 1, 2, 3 for A, B, start, select
}

func initJoypad() *Joypad {
	return &Joypad{
		dpad:    0x0,
		buttons: 0x0,
		sel:     0x30,
	}
}

func (j *Joypad) read() byte {
	var low byte = 0xf
	if !isBitSet(j.sel, 4) {
		low &= ^j.dpad
	}

	if !isBitSet(j.sel, 5) {
		low &= ^j.buttons
	}

	return 0xc0 | j.sel | low
}

func (j *Joypad) write(val byte) {
	// only bits 4 and 5, which we represent in sel, are writable
	j.sel = val & 0x30
}

// Update determines whether a button has been pressed this frame, determining whether
// a joypad interrupt needs to be requested
func (j *Joypad) update(inp byte) bool {
	// if we derive that anything has been pressed, we return a JoypadInterruptRequest
	oldJoypad := j.read()
	j.buttons = (inp & 0x0f)
	j.dpad = (inp >> 4)
	newJoypad := j.read()

	return oldJoypad&^newJoypad&0x0f != 0
}
