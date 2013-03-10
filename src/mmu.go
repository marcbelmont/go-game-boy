package main

import "io/ioutil"

////////////
// Memory //
////////////

type MMU struct {
	inbios, mif, mie            uint16
	bios, rom, wram, eram, zram []byte
	romoffs, ramoffs            uint16
	carttype                    byte
	mbc                         []MBC
}

type MBC struct {
	rombank, rambank, ramon, mode uint16
}

func (m *MMU) reset() {
	m.wram = make([]byte, 8192)
	m.eram = make([]byte, 32768)
	m.zram = make([]byte, 127)
	m.inbios = 0
	m.mie = 0
	m.mif = 0
	m.carttype = 0
	m.romoffs = 0x4000
	m.ramoffs = 0
	m.mbc = make([]MBC, 2)
	m.mbc[1] = MBC{}
}

func (m *MMU) load(filename string) {
	buf, err := ioutil.ReadFile(filename)
	if err == nil {
		m.rom = buf
		m.carttype = mmu.rom[0x0147]
	} else {
		panic(err)
	}
}

// Read write

func (m *MMU) rb(addr uint16) uint16 {
	switch addr & 0xF000 {
	case 0x0000:
		if m.inbios == 1 {
			if addr < 0x0100 {
				return uint16(m.bios[addr])
			}
			if z80.pc == 0x0100 {
				m.inbios = 0
			}
		} else {
			return uint16(m.rom[addr])
		}

	case 0x1000, 0x2000, 0x3000:
		return uint16(m.rom[addr])

		// ROM bank 1
	case 0x4000, 0x5000, 0x6000, 0x7000:
		return uint16(m.rom[m.romoffs+(addr&0x3fff)])

		// VRAM
	case 0x8000, 0x9000:
		return uint16(gpu.vram[addr&0x1FFF])

	// External RAM
	case 0xA000, 0xB000:
		return uint16(m.eram[addr&0x1FFF])

	// Work RAM and echo
	case 0xC000, 0xD000, 0xE000:
		return uint16(m.wram[addr&0x1FFF])

	// Everything else
	case 0xF000:
		switch addr & 0x0F00 {
		// Echo RAM
		case 0x000, 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xA00, 0xB00, 0xC00, 0xD00:
			return uint16(m.wram[addr&0x1FFF])

			// OAM
		case 0xE00:
			if (addr & 0xFF) < 0xA0 {
				return uint16(gpu.oam[addr&0xFF])
			} else {
				return 0
			}

			// Zeropage RAM, I/O
		case 0xF00:
			if addr == 0xFFFF {
				return uint16(m.mie)
			} else if addr > 0xFF7F {
				return uint16(m.zram[addr&0x7F])
			} else {
				switch addr & 0xF0 {
				case 0x00:
					switch addr & 0xF {
					case 0:
						return uint16(key.rb())
					case 4, 5, 6, 7:
						return uint16(timer.rb(addr))
					case 15:
						return uint16(m.mif) // Interrupt flags
					default:
						return 0
					}

				case 0x10, 0x20, 0x30:
					return 0

				case 0x40, 0x50, 0x60, 0x70:
					return uint16(gpu.rb(addr))
				}
			}
		}
	}
	return 0
}

func (m *MMU) rw(addr uint16) uint16 {
	return m.rb(addr) + (m.rb(addr+1) << 8)
}

func (m *MMU) wb(addr uint16, val uint16) {
	switch addr & 0xF000 {
	// ROM bank 0
	case 0x0000, 0x1000:
		if m.carttype == 1 {
			if (val & 0xF) == 0xA {
				m.mbc[1].ramon = 1
			} else {
				m.mbc[1].ramon = 0
			}
		}
	// fall through
	case 0x2000, 0x3000:
		if m.carttype == 1 {
			m.mbc[1].rombank &= 0x60
			val &= 0x1f
			if val == 0 {
				val = 1
			}
			m.mbc[1].rombank |= val
			m.romoffs = m.mbc[1].rombank * 0x4000
		}

	// ROM bank 1
	case 0x4000, 0x5000:
		if m.carttype == 1 {
			if m.mbc[1].mode != 0 {
				m.mbc[1].rombank = val & 3
				m.ramoffs = m.mbc[1].rambank * 0x2000
			} else {
				m.mbc[1].rombank &= 0x1f
				m.mbc[1].rombank |= ((val & 3) << 5)
				m.romoffs = m.mbc[1].rombank * 0x4000
			}
		}

	case 0x6000, 0x7000:
		if m.carttype == 1 {
			m.mbc[1].mode = val & 1
		}
	// VRAM
	case 0x8000, 0x9000:
		gpu.vram[addr&0x1FFF] = uint8(val)
		gpu.updatetile(addr&0x1FFF, uint8(val))
		break

	// External RAM
	case 0xA000, 0xB000:
		m.eram[m.ramoffs+addr&0x1FFF] = uint8(val)

	// Work RAM and echo
	case 0xC000, 0xD000, 0xE000:
		m.wram[addr&0x1FFF] = uint8(val)

	// Everything else
	case 0xF000:
		switch addr & 0x0F00 {
		// Echo RAM
		case 0x000, 0x100, 0x200, 0x300, 0x400, 0x500, 0x600, 0x700, 0x800, 0x900, 0xA00, 0xB00, 0xC00, 0xD00:
			m.wram[addr&0x1FFF] = uint8(val)

			// OAM
		case 0xE00:
			if (addr & 0xFF) < 0xA0 {
				gpu.oam[uint8(addr&0xFF)] = uint8(val)
			}
			gpu.updateoam(addr, uint8(val))

		// Zeropage RAM, I/O
		case 0xF00:
			if addr == 0xffff {
				m.mie = val
			} else if addr > 0xFF7F {
				m.zram[addr&0x7F] = uint8(val)
			} else {
				switch addr & 0xF0 {
				case 0:
					switch addr & 0xF {
					case 0:
						key.wb(byte(val))
					case 4, 5, 6, 7:
						timer.wb(addr, byte(val))
					case 15:
						m.mif = val
					}
				// case 0x10,0x20,0x30:
				// 1
				case 0x40, 0x50, 0x60, 0x70:
					gpu.wb(addr, byte(val))
				}
			}
		}
	}
}

func (m *MMU) ww(addr uint16, val uint16) {
	m.wb(addr, val&255)
	m.wb(addr+1, val>>8)
}

// ///////// //
// Constants //
// ///////// //

var BIOS = []byte{0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
	0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
	0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
	0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
	0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
	0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
	0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
	0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
	0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xF2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
	0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
	0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3c, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x4C,
	0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
	0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50}

const (
	HEIGHT, WIDTH = 144, 160
)
