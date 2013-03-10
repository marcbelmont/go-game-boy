package main

type TIMER struct {
	div, tma, tima, tac byte
	cmain, sub, cdiv    uint
}

func (t *TIMER) reset() {
	t.div = 0
	t.tma = 0
	t.tima = 0
	t.tac = 0
	t.cmain = 0
	t.sub = 0
	t.cdiv = 0
}
func (t *TIMER) step() {
	t.tima++
	t.cmain = 0
	if t.tima > 255 {
		t.tima = t.tma
		mmu.mif |= 4
	}
}
func (t *TIMER) inc() {
	t.sub += z80.m
	if t.sub > 3 {
		t.cmain++
		t.sub -= 4

		t.cdiv++
		if t.cdiv == 16 {
			t.cdiv = 0
			t.div++
			t.div &= 255
		}
	}

	if t.tac&4 != 0 {
		switch t.tac & 3 {
		case 0:
			if t.cmain >= 64 {
				t.step()
			}
			break
		case 1:
			if t.cmain >= 1 {
				t.step()
			}
			break
		case 2:
			if t.cmain >= 4 {
				t.step()
			}
			break
		case 3:
			if t.cmain >= 16 {
				t.step()
			}
			break
		}
	}

}
func (t *TIMER) rb(addr uint16) byte {
	switch addr {
	case 0xFF04:
		return t.div
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	}
	return 0
}

func (t *TIMER) wb(addr uint16, val byte) {
	switch addr {
	case 0xFF04:
		t.div = 0
		break
	case 0xFF05:
		t.tima = val
		break
	case 0xFF06:
		t.tma = val
		break
	case 0xFF07:
		t.tac = val & 7
		break
	}
}
