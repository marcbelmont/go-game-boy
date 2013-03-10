package main

import "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"

type KEY struct {
	keys   []byte
	colidx byte
}

func (k *KEY) reset() {
	k.keys = []byte{0x0f, 0x0f}
	k.colidx = 0
}

func (k *KEY) rb() byte {
	switch k.colidx {
	case 0x00:
		return 0x00
		break
	case 0x10:
		return k.keys[0]
		break
	case 0x20:
		return k.keys[1]
		break
	default:
		return 0x00
		break
	}
	return 0
}

func (k *KEY) wb(v byte) {
	k.colidx = v & 0x30
}

func (k *KEY) keydown(code uint32) {
	switch code {
	case sdl.K_RIGHT:
		k.keys[1] &= 0xE
		break
	case sdl.K_LEFT:
		k.keys[1] &= 0xD
		break
	case sdl.K_UP:
		k.keys[1] &= 0xB
		break
	case sdl.K_DOWN:
		k.keys[1] &= 0x7
		break
	case sdl.K_z:
		k.keys[0] &= 0xE
		break
	case sdl.K_x:
		k.keys[0] &= 0xD
		break
	case sdl.K_a:
		k.keys[0] &= 0xB
		break
	case sdl.K_r:
		k.keys[0] &= 0x7
		break
	}
}

func (k *KEY) keyup(code uint32) {
	switch code {
	case sdl.K_RIGHT:
		k.keys[1] |= 0x1
		break
	case sdl.K_LEFT:
		k.keys[1] |= 0x2
		break
	case sdl.K_UP:
		k.keys[1] |= 0x4
		break
	case sdl.K_DOWN:
		k.keys[1] |= 0x8
		break
	case sdl.K_z:
		k.keys[0] |= 0x1
		break
	case sdl.K_x:
		k.keys[0] |= 0x2
		break
	case sdl.K_a:
		k.keys[0] |= 0x5
		break
	case sdl.K_r:
		k.keys[0] |= 0x8
		break
	}
}
