package main

import (
	"fmt"
)

type Z80 struct {
	// Time clocks
	cm, ct uint

	// Registers
	a, b, c, d, e, f, h, l         uint16
	ra, rb, rc, rd, re, rf, rh, rl uint16
	pc, sp, i, r, ime              uint16

	// Last instruction clock
	m, t uint

	// Instructions
	imap  []func(*Z80)
	cbmap []func(*Z80)

	// Others
	halt, stop uint
}

/////////////
// Helpers //
/////////////
var xxx = 0

func (z *Z80) exec() {
	if z.halt == 1 {
		z.m = 1
	} else {
		x := z.pc
		z.pc++
		addr := mmu.rb(x)
		z.imap[addr](z)
		z.pc &= 65535
	}
	if z.ime != 0 && mmu.mie != 0 && mmu.mif != 0 {
		z.halt = 0
		z.ime = 0
		ifired := mmu.mie & mmu.mif
		if ifired&1 != 0 {
			mmu.mif &= 0xFE
			z.RST40()
		} else if ifired&2 != 0 {
			mmu.mif &= 0xFD
			z.RST48()
		} else if ifired&4 != 0 {
			mmu.mif &= 0xFB
			z.RST50()
		} else if ifired&8 != 0 {
			mmu.mif &= 0xF7
			z.RST58()
		} else if ifired&16 != 0 {
			mmu.mif &= 0xEF
			z.RST60()
		} else {
			z.ime = 1
		}

	}
	z.cm += z.m
	z.ct += z.t
}

func (z *Z80) print() {
	fmt.Printf(" a:%4X  b:%4X c:%4X d:%4X e:%4X f:%4X h:%4X l:%4X\n",
		z.a, z.b, z.c, z.d, z.e, z.f, z.h, z.l)
	fmt.Printf("pc:%4X sp:%4X\n", z.pc, z.sp)
	fmt.Printf(" m:%4X  t:%4X cm:%4X ct:%4X\n", z.m, z.t, z.cm, z.ct)
}

func (z *Z80) printi() {
	fmt.Printf(" a:%4v  b:%4v c:%4v d:%4v e:%4v f:%4v h:%4v l:%4v\n",
		z.a, z.b, z.c, z.d, z.e, z.f, z.h, z.l)
	fmt.Printf("pc:%4v sp:%4v\n", z.pc, z.sp)
	fmt.Printf(" m:%4v  t:%4v cm:%4v ct:%4v\n", z.m, z.t, z.cm, z.ct)
}

func (z *Z80) reset() {
	z.a, z.b, z.c, z.d, z.e, z.f, z.h, z.l = 0, 0, 0, 0, 0, 0, 0, 0
	z.pc, z.sp = 0, 0
	z.cm, z.ct = 0, 0
	z.halt, z.stop = 0, 0
	z.i, z.r, z.ime = 0, 0, 1

	z80.pc = 0x100
	z80.sp = 0xFFFE
	z80.c = 0x13
	z80.e = 0xD8
	z80.a = 1
}

/////////
// ASM //
/////////

// Load / store

func (z *Z80) LDrr_bb() { z.b = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_bc() { z.b = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_bd() { z.b = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_be() { z.b = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_bh() { z.b = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_bl() { z.b = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ba() { z.b = z.a; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_cb() { z.c = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_cc() { z.c = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_cd() { z.c = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ce() { z.c = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ch() { z.c = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_cl() { z.c = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ca() { z.c = z.a; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_db() { z.d = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_dc() { z.d = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_dd() { z.d = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_de() { z.d = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_dh() { z.d = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_dl() { z.d = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_da() { z.d = z.a; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_eb() { z.e = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ec() { z.e = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ed() { z.e = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ee() { z.e = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_eh() { z.e = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_el() { z.e = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ea() { z.e = z.a; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_hb() { z.h = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_hc() { z.h = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_hd() { z.h = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_he() { z.h = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_hh() { z.h = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_hl() { z.h = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ha() { z.h = z.a; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_lb() { z.l = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_lc() { z.l = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ld() { z.l = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_le() { z.l = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_lh() { z.l = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ll() { z.l = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_la() { z.l = z.a; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ab() { z.a = z.b; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ac() { z.a = z.c; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ad() { z.a = z.d; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ae() { z.a = z.e; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_ah() { z.a = z.h; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_al() { z.a = z.l; z.m = 1; z.t = 4 }
func (z *Z80) LDrr_aa() { z.a = z.a; z.m = 1; z.t = 4 }

func (z *Z80) LDrHLm_b() { z.b = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDrHLm_c() { z.c = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDrHLm_d() { z.d = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDrHLm_e() { z.e = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDrHLm_h() { z.h = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDrHLm_l() { z.l = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDrHLm_a() { z.a = mmu.rb((uint16(z.h) << 8) + uint16(z.l)); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_b() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.b); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_c() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.c); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_d() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.d); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_e() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.e); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_h() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.h); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_l() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.l); z.m = 2; z.t = 8 }
func (z *Z80) LDHLmr_a() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.a); z.m = 2; z.t = 8 }
func (z *Z80) LDrn_b()   { z.b = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }
func (z *Z80) LDrn_c()   { z.c = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }
func (z *Z80) LDrn_d()   { z.d = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }
func (z *Z80) LDrn_e()   { z.e = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }
func (z *Z80) LDrn_h()   { z.h = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }
func (z *Z80) LDrn_l()   { z.l = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }
func (z *Z80) LDrn_a()   { z.a = mmu.rb(z.pc); z.pc++; z.m = 2; z.t = 8 }

func (z *Z80) LDHLmn() { mmu.wb((uint16(z.h)<<8)+uint16(z.l), mmu.rb(z.pc)); z.pc++; z.m = 3; z.t = 12 }

func (z *Z80) LDBCmA() { mmu.wb((uint16(z.b)<<8)+uint16(z.c), z.a); z.m = 2; z.t = 8 }
func (z *Z80) LDDEmA() { mmu.wb((uint16(z.d)<<8)+uint16(z.e), z.a); z.m = 2; z.t = 8 }

func (z *Z80) LDmmA() { mmu.wb(mmu.rw(z.pc), z.a); z.pc += 2; z.m = 4; z.t = 16 }

func (z *Z80) LDABCm() {
	z.a = mmu.rb((uint16(z.b) << 8) + uint16(z.c))
	z.m = 2
	z.t = 8
}
func (z *Z80) LDADEm() { z.a = mmu.rb((uint16(z.d) << 8) + uint16(z.e)); z.m = 2; z.t = 8 }

func (z *Z80) LDAmm() { z.a = mmu.rb(mmu.rw(z.pc)); z.pc += 2; z.m = 4; z.t = 16 }

func (z *Z80) LDBCnn() { z.c = mmu.rb(z.pc); z.b = mmu.rb(z.pc + 1); z.pc += 2; z.m = 3; z.t = 12 }
func (z *Z80) LDDEnn() { z.e = mmu.rb(z.pc); z.d = mmu.rb(z.pc + 1); z.pc += 2; z.m = 3; z.t = 12 }
func (z *Z80) LDHLnn() { z.l = mmu.rb(z.pc); z.h = mmu.rb(z.pc + 1); z.pc += 2; z.m = 3; z.t = 12 }
func (z *Z80) LDSPnn() { z.sp = mmu.rw(z.pc); z.pc += 2; z.m = 3; z.t = 12 }

func (z *Z80) LDHLmm() {
	var i = mmu.rw(z.pc)
	z.pc += 2
	z.l = mmu.rb(i)
	z.h = mmu.rb(i + 1)
	z.m = 5
	z.t = 20
}
func (z *Z80) LDmmHL() {
	var i = mmu.rw(z.pc)
	z.pc += 2
	mmu.ww(i, (uint16(z.h)<<8)+uint16(z.l))
	z.m = 5
	z.t = 20
}

func (z *Z80) LDHLIA() {
	mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.a)
	z.l = (z.l + 1) & 255
	if z.l == 0 {
		z.h = (z.h + 1) & 255
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) LDAHLI() {
	z.a = mmu.rb((uint16(z.h) << 8) + uint16(z.l))
	z.l = (z.l + 1) & 255
	if z.l == 0 {
		z.h = (z.h + 1) & 255
	}
	z.m = 2
	z.t = 8
}

func (z *Z80) LDHLDA() {
	mmu.wb((uint16(z.h)<<8)+uint16(z.l), z.a)
	z.l = (z.l - 1) & 255
	if z.l == 255 {
		z.h = (z.h - 1) & 255
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) LDAHLD() {
	z.a = mmu.rb((uint16(z.h) << 8) + uint16(z.l))
	z.l = (z.l - 1) & 255
	if z.l == 255 {
		z.h = (z.h - 1) & 255
	}
	z.m = 2
	z.t = 8
}

func (z *Z80) LDAIOn() { z.a = mmu.rb(0xFF00 + uint16(mmu.rb(z.pc))); z.pc++; z.m = 3; z.t = 12 }
func (z *Z80) LDIOnA() {
	mmu.wb(0xFF00+uint16(mmu.rb(z.pc)), z.a)
	z.pc++
	z.m = 3
	z.t = 12
}
func (z *Z80) LDAIOC() { z.a = mmu.rb(0xFF00 + uint16(z.c)); z.m = 2; z.t = 8 }
func (z *Z80) LDIOCA() { mmu.wb(0xFF00+uint16(z.c), z.a); z.m = 2; z.t = 8 }

func (z *Z80) LDHLSPn() {
	var i = uint16(mmu.rb(z.pc))
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	i += z.sp
	z.h = uint16(i>>8) & 255
	z.l = uint16(i & 255)
	z.m = 3
	z.t = 12
}

func (z *Z80) SWAPr_b() {
	var tr = z.b
	z.b = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.b {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}
func (z *Z80) SWAPr_c() {
	var tr = z.c
	z.c = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.c {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}
func (z *Z80) SWAPr_d() {
	var tr = z.d
	z.d = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.d {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}
func (z *Z80) SWAPr_e() {
	var tr = z.e
	z.e = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.e {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}
func (z *Z80) SWAPr_h() {
	var tr = z.h
	z.h = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.h {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}
func (z *Z80) SWAPr_l() {
	var tr = z.l
	z.l = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.l {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}
func (z *Z80) SWAPr_a() {
	var tr = z.a
	z.a = ((tr & 0xF) << 4) | ((tr & 0xF0) >> 4)
	if 0 != z.a {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
}

/*--- Data processing ---*/
func (z *Z80) ADDr_b() {
	var a = z.a
	z.a += z.b
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.b^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDr_c() {
	var a = z.a
	z.a += z.c
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.c^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDr_d() {
	var a = z.a
	z.a += z.d
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.d^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDr_e() {
	var a = z.a
	z.a += z.e
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.e^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDr_h() {
	var a = z.a
	z.a += z.h
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.h^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDr_l() {
	var a = z.a
	z.a += z.l
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.l^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDr_a() {
	var a = z.a
	z.a += z.a
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.a^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) ADDHL() {
	var a = z.a
	var m = mmu.rb((z.h << 8) + z.l)
	z.a += m
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^a^m)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 2
}
func (z *Z80) ADDn() {
	var a = z.a
	var m = mmu.rb(z.pc)
	z.a += m
	z.pc++
	if z.a > 255 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^a^m)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 2
}

func (z *Z80) ADDHLBC() {
	var hl = (uint16(z.h) << 8) + uint16(z.l)
	hl += (uint16(z.b) << 8) + uint16(z.c)
	if hl > 65535 {
		z.f |= 0x10
	} else {
		z.f &= 0xEF
	}
	z.h = uint16((hl >> 8) & 255)
	z.l = uint16(hl & 255)
	z.m = 3
	z.t = 12
}
func (z *Z80) ADDHLDE() {
	var hl = (uint16(z.h) << 8) + uint16(z.l)
	hl += (uint16(z.d) << 8) + uint16(z.e)
	if hl > 65535 {
		z.f |= 0x10
	} else {
		z.f &= 0xEF
	}
	z.h = uint16((hl >> 8) & 255)
	z.l = uint16(hl & 255)
	z.m = 3
	z.t = 12
}
func (z *Z80) ADDHLHL() {
	var hl = (uint16(z.h) << 8) + uint16(z.l)
	hl += (uint16(z.h) << 8) + uint16(z.l)
	if hl > 65535 {
		z.f |= 0x10
	} else {
		z.f &= 0xEF
	}
	z.h = uint16((hl >> 8) & 255)
	z.l = uint16(hl & 255)
	z.m = 3
	z.t = 12
}
func (z *Z80) ADDHLSP() {
	var hl = (uint16(z.h) << 8) + uint16(z.l)
	hl += z.sp
	if hl > 65535 {
		z.f |= 0x10
	} else {
		z.f &= 0xEF
	}
	z.h = uint16((hl >> 8) & 255)
	z.l = uint16(hl & 255)
	z.m = 3
	z.t = 12
}
func (z *Z80) ADDSPn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.sp += uint16(i)
	z.m = 4
	z.t = 16
}

func (z *Z80) ADCr_b() {
	var a = z.a
	z.a += z.b
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.b ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCr_c() {
	var a = z.a
	z.a += z.c
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.c ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCr_d() {
	var a = z.a
	z.a += z.d
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.d ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCr_e() {
	var a = z.a
	z.a += z.e
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.e ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCr_h() {
	var a = z.a
	z.a += z.h
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.h ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCr_l() {
	var a = z.a
	z.a += z.l
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.l ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCr_a() {
	var a = z.a
	z.a += z.a
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f |= 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.a ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) ADCHL() {
	var a = z.a
	var m = mmu.rb((z.h << 8) + z.l)
	z.a += m
	if z.f&0x10 != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f = 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^m^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 2
}
func (z *Z80) ADCn() {
	var a = z.a
	var m = mmu.rb(z.pc)
	z.a += m
	z.pc++
	if (z.f & 0x10) != 0 {
		z.a += 1
	}
	if z.a > 255 {
		z.f = 0x10
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ m ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 2
}

func (z *Z80) SUBr_b() {
	var a = z.a
	if z.a < z.b {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.b) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.b ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) SUBr_c() {
	var a = z.a
	if z.a < z.c {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.c) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.c ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) SUBr_d() {
	var a = z.a
	if z.a < z.d {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.d) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.d ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) SUBr_e() {
	var a = z.a
	if z.a < z.e {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.e) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.e ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) SUBr_h() {
	var a = z.a
	if z.a < z.h {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.h) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.h ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4

}
func (z *Z80) SUBr_l() {
	var a = z.a
	if z.a < z.l {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.l) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.l ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) SUBr_a() {
	var a = z.a
	if z.a < z.a {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - z.a) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.a ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) SUBHL() {
	var a = z.a
	var m = mmu.rb((z.h << 8) + z.l)
	if z.a < m {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - m) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^m^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 2
}
func (z *Z80) SUBn() {
	var a = z.a
	var m = mmu.rb(z.pc)
	z.pc++
	if z.a < m {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a = (z.a - m) & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^m^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 2
}

func (z *Z80) SBCr_b() {
	var a = z.a
	z.a -= z.b
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.b^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCr_c() {
	var a = z.a
	z.a -= z.c
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.c^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCr_d() {
	var a = z.a
	z.a -= z.d
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.d^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCr_e() {
	var a = z.a
	z.a -= z.e
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.e^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCr_h() {
	var a = z.a
	z.a -= z.h
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.h^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCr_l() {
	var a = z.a
	z.a -= z.l
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.l^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCr_a() {
	var a = z.a
	z.a -= z.a
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if (z.a^z.a^a)&0x10 != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) SBCHL() {
	var a = z.a
	var m = mmu.rb((z.h << 8) + z.l)
	z.a -= m
	if z.f&0x10 != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ m ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 2
}
func (z *Z80) SBCn() {
	var a = z.a
	var m = mmu.rb(z.pc)
	z.a -= m
	z.pc++
	if (z.f & 0x10) != 0 {
		z.a -= 1
	}
	if z.a > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}

	z.a &= 255
	if z.a == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ m ^ a) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 2
}

func (z *Z80) CPr_b() {
	var i = z.a
	i -= z.b
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.b ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPr_c() {
	var i = z.a
	i -= z.c
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.c ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPr_d() {
	var i = z.a
	i -= z.d
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.d ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPr_e() {
	var i = z.a
	i -= z.e
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.e ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPr_h() {
	var i = z.a
	i -= z.h
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.h ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPr_l() {
	var i = z.a
	i -= z.l
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.l ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPr_a() {
	var i = z.a
	i -= z.a
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ z.a ^ i) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 1
}
func (z *Z80) CPHL() {
	var i = z.a
	var m = mmu.rb((z.h << 8) + z.l)
	i -= m
	if i > 40000 {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i &= 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ i ^ m) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 2
}
func (z *Z80) CPn() {
	var i = z.a
	var m = mmu.rb(z.pc)
	z.pc++
	if i < m {
		z.f = 0x50
	} else {
		z.f = 0x40
	}
	i = (i - m) & 255
	if i == 0 {
		z.f |= 0x80
	}
	if ((z.a ^ i ^ m) & 0x10) != 0 {
		z.f |= 0x20
	}
	z.m = 2
}

func (z *Z80) DAA() {
	var a = z.a
	if (z.f&0x20) != 0 || ((z.a & 15) > 9) {
		z.a += 6
	}
	z.f &= 0xEF
	if (z.f&0x20) != 0 || (a > 0x99) {
		z.a += 0x60
		z.f |= 0x10
	}
	z.m = 1
}

func (z *Z80) ANDr_b() {
	z.a &= z.b
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDr_c() {
	z.a &= z.c
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDr_d() {
	z.a &= z.d
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDr_e() {
	z.a &= z.e
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDr_h() {
	z.a &= z.h
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDr_l() {
	z.a &= z.l
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDr_a() {
	z.a &= z.a
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ANDHL() {
	z.a &= mmu.rb((z.h << 8) + z.l)
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) ANDn() {
	z.a &= mmu.rb(z.pc)
	z.pc++
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}

func (z *Z80) ORr_b() {
	z.a |= z.b
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORr_c() {
	z.a |= z.c
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORr_d() {
	z.a |= z.d
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORr_e() {
	z.a |= z.e
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORr_h() {
	z.a |= z.h
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORr_l() {
	z.a |= z.l
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORr_a() {
	z.a |= z.a
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) ORHL() {
	z.a |= mmu.rb((z.h << 8) + z.l)
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) ORn() {
	z.a |= mmu.rb(z.pc)
	z.pc++
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}

func (z *Z80) XORr_b() {
	z.a ^= z.b
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORr_c() {
	z.a ^= z.c
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORr_d() {
	z.a ^= z.d
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORr_e() {
	z.a ^= z.e
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORr_h() {
	z.a ^= z.h
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORr_l() {
	z.a ^= z.l
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORr_a() {
	z.a ^= z.a
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) XORHL() {
	z.a ^= mmu.rb((z.h << 8) + z.l)
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) XORn() {
	z.a ^= mmu.rb(z.pc)
	z.pc++
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}

func (z *Z80) INCr_b() {
	z.b++
	z.b &= 255
	if z.b == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCr_c() {
	z.c++
	z.c &= 255
	if z.c == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCr_d() {
	z.d++
	z.d &= 255
	if z.d == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCr_e() {
	z.e++
	z.e &= 255
	if z.e == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCr_h() {
	z.h++
	z.h &= 255
	if z.h == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCr_l() {
	z.l++
	z.l &= 255
	if z.l == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCr_a() {
	z.a++
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) INCHLm() {
	var i = mmu.rb((z.h<<8)+z.l) + 1
	i &= 255
	mmu.wb((z.h<<8)+z.l, i)
	if i == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) DECr_b() {
	z.b--
	z.b &= 255
	if z.b == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECr_c() {
	z.c--
	z.c &= 255
	if z.c == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECr_d() {
	z.d--
	z.d &= 255
	if z.d == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECr_e() {
	z.e--
	z.e &= 255
	if z.e == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECr_h() {
	z.h--
	z.h &= 255
	if z.h == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECr_l() {
	z.l--
	z.l &= 255
	if z.l == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECr_a() {
	z.a--
	z.a &= 255
	if z.a == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 1
	z.t = 4
}
func (z *Z80) DECHLm() {
	var i = mmu.rb((z.h<<8)+z.l) - 1
	i &= 255
	mmu.wb((z.h<<8)+z.l, i)
	if i == 0 {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 3
	z.t = 12
}

func (z *Z80) INCBC() {
	z.c = (z.c + 1) & 255
	if z.c == 0 {
		z.b = (z.b + 1) & 255
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) INCDE() {
	z.e = (z.e + 1) & 255
	if z.e == 0 {
		z.d = (z.d + 1) & 255
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) INCHL() {
	z.l = (z.l + 1) & 255
	if z.l == 0 {
		z.h = (z.h + 1) & 255
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) INCSP() { z.sp = (z.sp + 1) & 65535; z.m = 1; z.t = 4 }

func (z *Z80) DECBC() {
	z.c = (z.c - 1) & 255
	if z.c == 255 {
		z.b = (z.b - 1) & 255
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) DECDE() {
	z.e = (z.e - 1) & 255
	if z.e == 255 {
		z.d = (z.d - 1) & 255
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) DECHL() {
	z.l = (z.l - 1) & 255
	if z.l == 255 {
		z.h = (z.h - 1) & 255
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) DECSP() { z.sp = (z.sp - 1) & 65535; z.m = 1; z.t = 4 }

/*--- Bit manipulation ---*/
func (z *Z80) BIT0b() {
	if 0 == (z.b & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0c() {
	if 0 == (z.c & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0d() {
	if 0 == (z.d & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0e() {
	if 0 == (z.e & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0h() {
	if 0 == (z.h & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0l() {
	if 0 == (z.l & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0a() {
	if 0 == (z.a & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT0m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x01) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 3
	z.t = 12
}

func (z *Z80) RES0b() { z.b &= 0xFE; z.m = 2 }
func (z *Z80) RES0c() { z.c &= 0xFE; z.m = 2 }
func (z *Z80) RES0d() { z.d &= 0xFE; z.m = 2 }
func (z *Z80) RES0e() { z.e &= 0xFE; z.m = 2 }
func (z *Z80) RES0h() { z.h &= 0xFE; z.m = 2 }
func (z *Z80) RES0l() { z.l &= 0xFE; z.m = 2 }
func (z *Z80) RES0a() { z.a &= 0xFE; z.m = 2 }
func (z *Z80) RES0m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xFE
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET0b() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0c() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0d() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0e() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0h() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0l() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0a() { z.b |= 0x01; z.m = 2 }
func (z *Z80) SET0m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x01
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}
func (z *Z80) BIT1b() {
	if 0 == (z.b & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1c() {
	if 0 == (z.c & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1d() {
	if 0 == (z.d & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1e() {
	if 0 == (z.e & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1h() {
	if 0 == (z.h & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1l() {
	if 0 == (z.l & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1a() {
	if 0 == (z.a & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 2
	z.t = 8
}
func (z *Z80) BIT1m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x02) {
		z.f = 0x80
	} else {
		z.f = 0
	}

	z.m = 3
	z.t = 12
}

func (z *Z80) RES1b() { z.b &= 0xFD; z.m = 2 }
func (z *Z80) RES1c() { z.c &= 0xFD; z.m = 2 }
func (z *Z80) RES1d() { z.d &= 0xFD; z.m = 2 }
func (z *Z80) RES1e() { z.e &= 0xFD; z.m = 2 }
func (z *Z80) RES1h() { z.h &= 0xFD; z.m = 2 }
func (z *Z80) RES1l() { z.l &= 0xFD; z.m = 2 }
func (z *Z80) RES1a() { z.a &= 0xFD; z.m = 2 }
func (z *Z80) RES1m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xFD
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET1b() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1c() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1d() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1e() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1h() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1l() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1a() { z.b |= 0x02; z.m = 2 }
func (z *Z80) SET1m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x02
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) RES2b() { z.b &= 0xFB; z.m = 2 }
func (z *Z80) RES2c() { z.c &= 0xFB; z.m = 2 }
func (z *Z80) RES2d() { z.d &= 0xFB; z.m = 2 }
func (z *Z80) RES2e() { z.e &= 0xFB; z.m = 2 }
func (z *Z80) RES2h() { z.h &= 0xFB; z.m = 2 }
func (z *Z80) RES2l() { z.l &= 0xFB; z.m = 2 }
func (z *Z80) RES2a() { z.a &= 0xFB; z.m = 2 }
func (z *Z80) RES2m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xFB
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET2b() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2c() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2d() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2e() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2h() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2l() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2a() { z.b |= 0x04; z.m = 2 }
func (z *Z80) SET2m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x04
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) RES3b() { z.b &= 0xF7; z.m = 2 }
func (z *Z80) RES3c() { z.c &= 0xF7; z.m = 2 }
func (z *Z80) RES3d() { z.d &= 0xF7; z.m = 2 }
func (z *Z80) RES3e() { z.e &= 0xF7; z.m = 2 }
func (z *Z80) RES3h() { z.h &= 0xF7; z.m = 2 }
func (z *Z80) RES3l() { z.l &= 0xF7; z.m = 2 }
func (z *Z80) RES3a() { z.a &= 0xF7; z.m = 2 }
func (z *Z80) RES3m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xF7
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET3b() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3c() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3d() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3e() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3h() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3l() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3a() { z.b |= 0x08; z.m = 2 }
func (z *Z80) SET3m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x08
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) RES4b() { z.b &= 0xEF; z.m = 2 }
func (z *Z80) RES4c() { z.c &= 0xEF; z.m = 2 }
func (z *Z80) RES4d() { z.d &= 0xEF; z.m = 2 }
func (z *Z80) RES4e() { z.e &= 0xEF; z.m = 2 }
func (z *Z80) RES4h() { z.h &= 0xEF; z.m = 2 }
func (z *Z80) RES4l() { z.l &= 0xEF; z.m = 2 }
func (z *Z80) RES4a() { z.a &= 0xEF; z.m = 2 }
func (z *Z80) RES4m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xEF
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET4b() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4c() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4d() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4e() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4h() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4l() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4a() { z.b |= 0x10; z.m = 2 }
func (z *Z80) SET4m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x10
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) RES5b() { z.b &= 0xDF; z.m = 2 }
func (z *Z80) RES5c() { z.c &= 0xDF; z.m = 2 }
func (z *Z80) RES5d() { z.d &= 0xDF; z.m = 2 }
func (z *Z80) RES5e() { z.e &= 0xDF; z.m = 2 }
func (z *Z80) RES5h() { z.h &= 0xDF; z.m = 2 }
func (z *Z80) RES5l() { z.l &= 0xDF; z.m = 2 }
func (z *Z80) RES5a() { z.a &= 0xDF; z.m = 2 }
func (z *Z80) RES5m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xDF
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET5b() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5c() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5d() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5e() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5h() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5l() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5a() { z.b |= 0x20; z.m = 2 }
func (z *Z80) SET5m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x20
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) RES6b() { z.b &= 0xBF; z.m = 2 }
func (z *Z80) RES6c() { z.c &= 0xBF; z.m = 2 }
func (z *Z80) RES6d() { z.d &= 0xBF; z.m = 2 }
func (z *Z80) RES6e() { z.e &= 0xBF; z.m = 2 }
func (z *Z80) RES6h() { z.h &= 0xBF; z.m = 2 }
func (z *Z80) RES6l() { z.l &= 0xBF; z.m = 2 }
func (z *Z80) RES6a() { z.a &= 0xBF; z.m = 2 }
func (z *Z80) RES6m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0xBF
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET6b() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6c() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6d() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6e() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6h() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6l() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6a() { z.b |= 0x40; z.m = 2 }
func (z *Z80) SET6m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x40
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) RES7b() { z.b &= 0x7F; z.m = 2 }
func (z *Z80) RES7c() { z.c &= 0x7F; z.m = 2 }
func (z *Z80) RES7d() { z.d &= 0x7F; z.m = 2 }
func (z *Z80) RES7e() { z.e &= 0x7F; z.m = 2 }
func (z *Z80) RES7h() { z.h &= 0x7F; z.m = 2 }
func (z *Z80) RES7l() { z.l &= 0x7F; z.m = 2 }
func (z *Z80) RES7a() { z.a &= 0x7F; z.m = 2 }
func (z *Z80) RES7m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i &= 0x7F
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) SET7b() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7c() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7d() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7e() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7h() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7l() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7a() { z.b |= 0x80; z.m = 2 }
func (z *Z80) SET7m() {
	var i = mmu.rb((z.h << 8) + z.l)
	i |= 0x80
	mmu.wb((z.h<<8)+z.l, i)
	z.m = 4
}

func (z *Z80) BIT2b() {
	if 0 == (z.b & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2c() {
	if 0 == (z.c & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2d() {
	if 0 == (z.d & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2e() {
	if 0 == (z.e & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2h() {
	if 0 == (z.h & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2l() {
	if 0 == (z.l & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2a() {
	if 0 == (z.a & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT2m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x04) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) BIT3b() {
	if 0 == (z.b & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3c() {
	if 0 == (z.c & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3d() {
	if 0 == (z.d & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3e() {
	if 0 == (z.e & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3h() {
	if 0 == (z.h & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3l() {
	if 0 == (z.l & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3a() {
	if 0 == (z.a & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT3m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x08) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) BIT4b() {
	if 0 == (z.b & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4c() {
	if 0 == (z.c & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4d() {
	if 0 == (z.d & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4e() {
	if 0 == (z.e & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4h() {
	if 0 == (z.h & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4l() {
	if 0 == (z.l & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4a() {
	if 0 == (z.a & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT4m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x10) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) BIT5b() {
	if 0 == (z.b & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5c() {
	if 0 == (z.c & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5d() {
	if 0 == (z.d & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5e() {
	if 0 == (z.e & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5h() {
	if 0 == (z.h & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5l() {
	if 0 == (z.l & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5a() {
	if 0 == (z.a & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT5m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x20) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) BIT6b() {
	if 0 == (z.b & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6c() {
	if 0 == (z.c & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6d() {
	if 0 == (z.d & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6e() {
	if 0 == (z.e & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6h() {
	if 0 == (z.h & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6l() {
	if 0 == (z.l & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6a() {
	if 0 == (z.a & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT6m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x40) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) BIT7b() {
	if 0 == (z.b & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7c() {
	if 0 == (z.c & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7d() {
	if 0 == (z.d & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7e() {
	if 0 == (z.e & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7h() {
	if 0 == (z.h & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7l() {
	if 0 == (z.l & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7a() {
	if 0 == (z.a & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 2
	z.t = 8
}
func (z *Z80) BIT7m() {
	if 0 == (mmu.rb((uint16(z.h)<<8)+uint16(z.l)) & 0x80) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.m = 3
	z.t = 12
}

func (z *Z80) RLA() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.a&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a << 1) + ci
	z.a &= 255
	z.f = (z.f & 0xEF) + co
	z.m = 1
	z.t = 4
}
func (z *Z80) RLCA() {
	var ci, co uint16 = 0, 0
	if z.a&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.a&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a << 1) + ci
	z.a &= 255
	z.f = (z.f & 0xEF) + co
	z.m = 1
	z.t = 4
}
func (z *Z80) RRA() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.a&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a >> 1) + ci
	z.a &= 255
	z.f = (z.f & 0xEF) + co
	z.m = 1
	z.t = 4
}
func (z *Z80) RRCA() {
	var ci, co uint16 = 0, 0
	if z.a&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.a&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a >> 1) + ci
	z.a &= 255
	z.f = (z.f & 0xEF) + co
	z.m = 1
	z.t = 4
}

func (z *Z80) RLr_b() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.b&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b << 1) + ci
	z.b &= 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLr_c() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.c&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c << 1) + ci
	z.c &= 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLr_d() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.d&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d << 1) + ci
	z.d &= 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLr_e() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.e&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e << 1) + ci
	z.e &= 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLr_h() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h << 1) + ci
	z.h &= 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLr_l() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l << 1) + ci
	z.l &= 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLr_a() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.a&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a << 1) + ci
	z.a &= 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLHL() {
	var ci, co uint16 = 0, 0
	var i = mmu.rb((z.h << 8) + z.l)
	if z.f&0x10 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if i&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	i = (i << 1) + ci
	i &= 255
	if 0 == (i) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	mmu.wb((z.h<<8)+z.l, i)
	z.f = (z.f & 0xEF) + co
	z.m = 4
	z.t = 16
}

func (z *Z80) RLCr_b() {
	var ci, co uint16 = 0, 0
	if z.b&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.b&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b << 1) + ci
	z.b &= 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCr_c() {
	var ci, co uint16 = 0, 0
	if z.c&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.c&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c << 1) + ci
	z.c &= 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCr_d() {
	var ci, co uint16 = 0, 0
	if z.d&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.d&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d << 1) + ci
	z.d &= 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCr_e() {
	var ci, co uint16 = 0, 0
	if z.e&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.e&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e << 1) + ci
	z.e &= 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCr_h() {
	var ci, co uint16 = 0, 0
	if z.h&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h << 1) + ci
	z.h &= 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCr_l() {
	var ci, co uint16 = 0, 0
	if z.l&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l << 1) + ci
	z.l &= 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCr_a() {
	var ci, co uint16 = 0, 0
	if z.a&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if z.a&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a << 1) + ci
	z.a &= 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RLCHL() {
	var ci, co uint16 = 0, 0
	var i = mmu.rb((z.h << 8) + z.l)
	if i&0x80 != 0 {
		ci = 1
	} else {
		ci = 0
	}
	if i&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	i = (i << 1) + ci
	i &= 255
	if 0 == (i) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	mmu.wb((z.h<<8)+z.l, i)
	z.f = (z.f & 0xEF) + co
	z.m = 4
	z.t = 16
}

func (z *Z80) RRr_b() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.b&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b >> 1) + ci
	z.b &= 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRr_c() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.c&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c >> 1) + ci
	z.c &= 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRr_d() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.d&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d >> 1) + ci
	z.d &= 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRr_e() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.e&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e >> 1) + ci
	z.e &= 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRr_h() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h >> 1) + ci
	z.h &= 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRr_l() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l >> 1) + ci
	z.l &= 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRr_a() {
	var ci, co uint16 = 0, 0
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.a&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a >> 1) + ci
	z.a &= 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRHL() {
	var ci, co uint16 = 0, 0
	var i = mmu.rb((z.h << 8) + z.l)
	if z.f&0x10 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if i&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	i = (i >> 1) + ci
	i &= 255
	mmu.wb((z.h<<8)+z.l, i)
	if 0 == (i) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 4
	z.t = 16
}

func (z *Z80) RRCr_b() {
	var ci, co uint16 = 0, 0
	if z.b&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.b&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b >> 1) + ci
	z.b &= 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCr_c() {
	var ci, co uint16 = 0, 0
	if z.c&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.c&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c >> 1) + ci
	z.c &= 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCr_d() {
	var ci, co uint16 = 0, 0
	if z.d&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.d&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d >> 1) + ci
	z.d &= 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCr_e() {
	var ci, co uint16 = 0, 0
	if z.e&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.e&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e >> 1) + ci
	z.e &= 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCr_h() {
	var ci, co uint16 = 0, 0
	if z.l&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h >> 1) + ci
	z.h &= 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCr_l() {
	var ci, co uint16 = 0, 0
	if z.l&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l >> 1) + ci
	z.l &= 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCr_a() {
	var ci, co uint16 = 0, 0
	if z.a&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if z.a&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a >> 1) + ci
	z.a &= 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) RRCHL() {
	var ci, co uint16 = 0, 0
	var i = mmu.rb((z.h << 8) + z.l)
	if i&1 != 0 {
		ci = 0x80
	} else {
		ci = 0
	}
	if i&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	i = (i >> 1) + ci
	i &= 255
	mmu.wb((z.h<<8)+z.l, i)
	if 0 == (i) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 4
	z.t = 16
}

func (z *Z80) SLAr_b() {
	var co uint16 = 0
	if z.b&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b << 1) & 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLAr_c() {
	var co uint16 = 0
	if z.c&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c << 1) & 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLAr_d() {
	var co uint16 = 0
	if z.d&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d << 1) & 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLAr_e() {
	var co uint16 = 0
	if z.e&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e << 1) & 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLAr_h() {
	var co uint16 = 0
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h << 1) & 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLAr_l() {
	var co uint16 = 0
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l << 1) & 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLAr_a() {
	var co uint16 = 0
	if z.a&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a << 1) & 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}

func (z *Z80) SLLr_b() {
	var co uint16 = 0
	if z.b&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b<<1)&255 + 1
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLLr_c() {
	var co uint16 = 0
	if z.c&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c<<1)&255 + 1
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLLr_d() {
	var co uint16 = 0
	if z.d&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d<<1)&255 + 1
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLLr_e() {
	var co uint16 = 0
	if z.e&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e<<1)&255 + 1
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLLr_h() {
	var co uint16 = 0
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h<<1)&255 + 1
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLLr_l() {
	var co uint16 = 0
	if z.l&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l<<1)&255 + 1
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SLLr_a() {
	var co uint16 = 0
	if z.a&0x80 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a<<1)&255 + 1
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}

func (z *Z80) SRAr_b() {
	var ci, co uint16 = 0, 0
	ci = z.b & 0x80
	if z.b&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = ((z.b >> 1) + ci) & 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRAr_c() {
	var ci, co uint16 = 0, 0
	ci = z.c & 0x80
	if z.c&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = ((z.c >> 1) + ci) & 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRAr_d() {
	var ci, co uint16 = 0, 0
	ci = z.d & 0x80
	if z.d&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = ((z.d >> 1) + ci) & 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRAr_e() {
	var ci, co uint16 = 0, 0
	ci = z.e & 0x80
	if z.e&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = ((z.e >> 1) + ci) & 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRAr_h() {
	var ci, co uint16 = 0, 0
	ci = z.h & 0x80
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = ((z.h >> 1) + ci) & 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRAr_l() {
	var ci, co uint16 = 0, 0
	ci = z.l & 0x80
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = ((z.l >> 1) + ci) & 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRAr_a() {
	var ci, co uint16 = 0, 0
	ci = z.a & 0x80
	if z.a&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = ((z.a >> 1) + ci) & 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}

func (z *Z80) SRLr_b() {
	var co uint16 = 0
	if z.b&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.b = (z.b >> 1) & 255
	if 0 == (z.b) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRLr_c() {
	var co uint16 = 0
	if z.c&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.c = (z.c >> 1) & 255
	if 0 == (z.c) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRLr_d() {
	var co uint16 = 0
	if z.d&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.d = (z.d >> 1) & 255
	if 0 == (z.d) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRLr_e() {
	var co uint16 = 0
	if z.e&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.e = (z.e >> 1) & 255
	if 0 == (z.e) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRLr_h() {
	var co uint16 = 0
	if z.h&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.h = (z.h >> 1) & 255
	if 0 == (z.h) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRLr_l() {
	var co uint16 = 0
	if z.l&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.l = (z.l >> 1) & 255
	if 0 == (z.l) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}
func (z *Z80) SRLr_a() {
	var co uint16 = 0
	if z.a&1 != 0 {
		co = 0x10
	} else {
		co = 0
	}
	z.a = (z.a >> 1) & 255
	if 0 == (z.a) {
		z.f = 0x80
	} else {
		z.f = 0
	}
	z.f = (z.f & 0xEF) + co
	z.m = 2
	z.t = 8
}

func (z *Z80) CPL() {
	z.a = (^z.a) & 255
	if (z.a) != 0 {
		z.f = 0
	} else {
		z.f = 0x80
	}
	z.m = 1
	z.t = 4
}
func (z *Z80) NEG() {
	if z.a != 0 {
		z.f = 0x10
	} else {
		z.f = 0
	}
	z.a = -z.a & 255
	if z.a == 0 {
		z.f |= 0x80
	}
	z.m = 2
	z.t = 8
}

func (z *Z80) CCF() {
	var ci uint16 = 0
	if z.f&0x10 != 0 {
		ci = 0
	} else {
		ci = 0x10
	}
	z.f = (z.f & 0xEF) + ci
	z.m = 1
	z.t = 4
}
func (z *Z80) SCF() { z.f |= 0x10; z.m = 1; z.t = 4 }

/*--- Stack ---*/
func (z *Z80) PUSHBC() {
	z.sp--
	mmu.wb(z.sp, z.b)
	z.sp--
	mmu.wb(z.sp, z.c)
	z.m = 3
	z.t = 12
}
func (z *Z80) PUSHDE() {
	z.sp--
	mmu.wb(z.sp, z.d)
	z.sp--
	mmu.wb(z.sp, z.e)
	z.m = 3
	z.t = 12
}
func (z *Z80) PUSHHL() {
	z.sp--
	mmu.wb(z.sp, z.h)
	z.sp--
	mmu.wb(z.sp, z.l)
	z.m = 3
	z.t = 12
}
func (z *Z80) PUSHAF() {
	z.sp--
	mmu.wb(z.sp, z.a)
	z.sp--
	mmu.wb(z.sp, z.f)
	z.m = 3
	z.t = 12
}

func (z *Z80) POPBC() {
	z.c = mmu.rb(z.sp)
	z.sp++
	z.b = mmu.rb(z.sp)
	z.sp++
	z.m = 3
	z.t = 12
}
func (z *Z80) POPDE() {
	z.e = mmu.rb(z.sp)
	z.sp++
	z.d = mmu.rb(z.sp)
	z.sp++
	z.m = 3
	z.t = 12
}
func (z *Z80) POPHL() {
	z.l = mmu.rb(z.sp)
	z.sp++
	z.h = mmu.rb(z.sp)
	z.sp++
	z.m = 3
	z.t = 12
}
func (z *Z80) POPAF() {
	z.f = mmu.rb(z.sp)
	z.sp++
	z.a = mmu.rb(z.sp)
	z.sp++
	z.m = 3
	z.t = 12
}

/*--- Jump ---*/
func (z *Z80) JPnn() {
	z.pc = mmu.rw(z.pc)
	z.m = 3
	z.t = 12
}
func (z *Z80) JPHL() { z.pc = (z.h << 8) + z.l; z.m = 1; z.t = 4 } // FIXME??
func (z *Z80) JPNZnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x80) == 0x00 {
		z.pc = mmu.rw(z.pc)
		z.m++
		z.t += 4
	} else {
		z.pc += 2
	}
}
func (z *Z80) JPZnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x80) == 0x80 {
		z.pc = mmu.rw(z.pc)
		z.m++
		z.t += 4
	} else {
		z.pc += 2
	}
}
func (z *Z80) JPNCnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x10) == 0x00 {
		z.pc = mmu.rw(z.pc)
		z.m++
		z.t += 4
	} else {
		z.pc += 2
	}
}
func (z *Z80) JPCnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x10) == 0x10 {
		z.pc = mmu.rw(z.pc)
		z.m++
		z.t += 4
	} else {
		z.pc += 2
	}
}

func (z *Z80) JRn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.m = 2
	z.t = 8
	z.pc += i
	z.m++
	z.t += 4
}
func (z *Z80) JRNZn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.m = 2
	z.t = 8
	if (z.f & 0x80) == 0x00 {
		z.pc += i
		z.m++
		z.t += 4
	}
}
func (z *Z80) JRZn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.m = 2
	z.t = 8
	if (z.f & 0x80) == 0x80 {
		z.pc += i
		z.m++
		z.t += 4
	}
}
func (z *Z80) JRNCn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.m = 2
	z.t = 8
	if (z.f & 0x10) == 0x00 {
		z.pc += i
		z.m++
		z.t += 4
	}
}
func (z *Z80) JRCn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.m = 2
	z.t = 8
	if (z.f & 0x10) == 0x10 {
		z.pc += i
		z.m++
		z.t += 4
	}
}

func (z *Z80) DJNZn() {
	var i = mmu.rb(z.pc)
	if i > 127 {
		i = -((^i + 1) & 255)
	}
	z.pc++
	z.m = 2
	z.t = 8
	z.b--
	if z.b != 0 {
		z.pc += i
		z.m++
		z.t += 4
	}
}

func (z *Z80) CALLnn() {
	z.sp -= 2
	mmu.ww(z.sp, z.pc+2)
	z.pc = mmu.rw(z.pc)
	z.m = 5
	z.t = 20
}

func (z *Z80) CALLNZnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x80) == 0x00 {
		z.sp -= 2
		mmu.ww(z.sp, z.pc+2)
		z.pc = mmu.rw(z.pc)
		z.m += 2
		z.t += 8
	} else {
		z.pc += 2
	}
}
func (z *Z80) CALLZnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x80) == 0x80 {
		z.sp -= 2
		mmu.ww(z.sp, z.pc+2)
		z.pc = mmu.rw(z.pc)
		z.m += 2
		z.t += 8
	} else {
		z.pc += 2
	}
}
func (z *Z80) CALLNCnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x10) == 0x00 {
		z.sp -= 2
		mmu.ww(z.sp, z.pc+2)
		z.pc = mmu.rw(z.pc)
		z.m += 2
		z.t += 8
	} else {
		z.pc += 2
	}
}
func (z *Z80) CALLCnn() {
	z.m = 3
	z.t = 12
	if (z.f & 0x10) == 0x10 {
		z.sp -= 2
		mmu.ww(z.sp, z.pc+2)
		z.pc = mmu.rw(z.pc)
		z.m += 2
		z.t += 8
	} else {
		z.pc += 2
	}
}

func (z *Z80) RET() { z.pc = mmu.rw(z.sp); z.sp += 2; z.m = 3; z.t = 12 }
func (z *Z80) RETI() {
	z.ime = 1
	z.rrs()
	z.pc = mmu.rw(z.sp)
	z.sp += 2
	z.m = 3
	z.t = 12
}
func (z *Z80) RETNZ() {
	z.m = 1
	z.t = 4
	if (z.f & 0x80) == 0x00 {
		z.pc = mmu.rw(z.sp)
		z.sp += 2
		z.m += 2
		z.t += 8
	}
}
func (z *Z80) RETZ() {
	z.m = 1
	z.t = 4
	if (z.f & 0x80) == 0x80 {
		z.pc = mmu.rw(z.sp)
		z.sp += 2
		z.m += 2
		z.t += 8
	}
}
func (z *Z80) RETNC() {
	z.m = 1
	z.t = 4
	if (z.f & 0x10) == 0x00 {
		z.pc = mmu.rw(z.sp)
		z.sp += 2
		z.m += 2
		z.t += 8
	}
}
func (z *Z80) RETC() {
	z.m = 1
	z.t = 4
	if (z.f & 0x10) == 0x10 {
		z.pc = mmu.rw(z.sp)
		z.sp += 2
		z.m += 2
		z.t += 8
	}
}

func (z *Z80) RST00() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x00
	z.m = 3
	z.t = 12
}
func (z *Z80) RST08() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x08
	z.m = 3
	z.t = 12
}
func (z *Z80) RST10() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x10
	z.m = 3
	z.t = 12
}
func (z *Z80) RST18() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x18
	z.m = 3
	z.t = 12
}
func (z *Z80) RST20() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x20
	z.m = 3
	z.t = 12
}
func (z *Z80) RST28() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x28
	z.m = 3
	z.t = 12
}
func (z *Z80) RST30() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x30
	z.m = 3
	z.t = 12
}
func (z *Z80) RST38() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x38
	z.m = 3
	z.t = 12
}
func (z *Z80) RST40() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x40
	z.m = 3
	z.t = 12
}
func (z *Z80) RST48() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x48
	z.m = 3
	z.t = 12
}
func (z *Z80) RST50() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x50
	z.m = 3
	z.t = 12
}
func (z *Z80) RST58() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x58
	z.m = 3
	z.t = 12
}
func (z *Z80) RST60() {
	z.rsv()
	z.sp -= 2
	mmu.ww(z.sp, z.pc)
	z.pc = 0x60
	z.m = 3
	z.t = 12
}

func (z *Z80) NOP()  { z.m = 1; z.t = 4 }
func (z *Z80) HALT() { z.halt = 1; z.m = 1; z.t = 4 }

func (z *Z80) DI() { z.ime = 0; z.m = 1; z.t = 4 }
func (z *Z80) EI() { z.ime = 1; z.m = 1; z.t = 4 }

func (z *Z80) MAPcb() {
	var i = mmu.rb(z.pc)
	z.pc++
	z.pc &= 65535
	if int(i) < len(z.cbmap) {
		z.cbmap[i](z)
	}
}

func (z *Z80) rsv() {
	z.ra = z.a
	z.rb = z.b
	z.rc = z.c
	z.rd = z.d
	z.re = z.e
	z.rf = z.f
	z.rh = z.h
	z.rl = z.l
}

func (z *Z80) rrs() {
	z.a = z.ra
	z.b = z.rb
	z.c = z.rc
	z.d = z.rd
	z.e = z.re
	z.f = z.rf
	z.h = z.rh
	z.l = z.rl
}

// undefined
func (z *Z80) XX() {
	// var opc = z.pc - 1
	// z.print()
	// panic(opc)
	z.stop = 1
}

//////////////////////
// Instructions map //
//////////////////////

func Map() []func(*Z80) {
	return []func(*Z80){(*Z80).NOP,
		(*Z80).LDBCnn,
		(*Z80).LDBCmA,
		(*Z80).INCBC,
		(*Z80).INCr_b,
		(*Z80).DECr_b,
		(*Z80).LDrn_b,
		(*Z80).RLCA,
		(*Z80).XX, // (*Z80).LDmmSP, // FIXME
		(*Z80).ADDHLBC,
		(*Z80).LDABCm,
		(*Z80).DECBC,
		(*Z80).INCr_c,
		(*Z80).DECr_c,
		(*Z80).LDrn_c,
		(*Z80).RRCA,

		// 10
		(*Z80).DJNZn,
		(*Z80).LDDEnn,
		(*Z80).LDDEmA,
		(*Z80).INCDE,
		(*Z80).INCr_d,
		(*Z80).DECr_d,
		(*Z80).LDrn_d,
		(*Z80).RLA,
		(*Z80).JRn,
		(*Z80).ADDHLDE,
		(*Z80).LDADEm,
		(*Z80).DECDE,
		(*Z80).INCr_e,
		(*Z80).DECr_e,
		(*Z80).LDrn_e,
		(*Z80).RRA,

		// 20
		(*Z80).JRNZn,
		(*Z80).LDHLnn,
		(*Z80).LDHLIA,
		(*Z80).INCHL,
		(*Z80).INCr_h,
		(*Z80).DECr_h,
		(*Z80).LDrn_h,
		(*Z80).DAA,
		(*Z80).JRZn,
		(*Z80).ADDHLHL,
		(*Z80).LDAHLI,
		(*Z80).DECHL,
		(*Z80).INCr_l,
		(*Z80).DECr_l,
		(*Z80).LDrn_l,
		(*Z80).CPL,

		// 30
		(*Z80).JRNCn,
		(*Z80).LDSPnn,
		(*Z80).LDHLDA,
		(*Z80).INCSP,
		(*Z80).INCHLm,
		(*Z80).DECHLm,
		(*Z80).LDHLmn,
		(*Z80).SCF,
		(*Z80).JRCn,
		(*Z80).ADDHLSP,
		(*Z80).LDAHLD,
		(*Z80).DECSP,
		(*Z80).INCr_a,
		(*Z80).DECr_a,
		(*Z80).LDrn_a,
		(*Z80).CCF,

		// 40
		(*Z80).LDrr_bb,
		(*Z80).LDrr_bc,
		(*Z80).LDrr_bd,
		(*Z80).LDrr_be,
		(*Z80).LDrr_bh,
		(*Z80).LDrr_bl,
		(*Z80).LDrHLm_b,
		(*Z80).LDrr_ba,
		(*Z80).LDrr_cb,
		(*Z80).LDrr_cc,
		(*Z80).LDrr_cd,
		(*Z80).LDrr_ce,
		(*Z80).LDrr_ch,
		(*Z80).LDrr_cl,
		(*Z80).LDrHLm_c,
		(*Z80).LDrr_ca,

		// 50
		(*Z80).LDrr_db,
		(*Z80).LDrr_dc,
		(*Z80).LDrr_dd,
		(*Z80).LDrr_de,
		(*Z80).LDrr_dh,
		(*Z80).LDrr_dl,
		(*Z80).LDrHLm_d,
		(*Z80).LDrr_da,
		(*Z80).LDrr_eb,
		(*Z80).LDrr_ec,
		(*Z80).LDrr_ed,
		(*Z80).LDrr_ee,
		(*Z80).LDrr_eh,
		(*Z80).LDrr_el,
		(*Z80).LDrHLm_e,
		(*Z80).LDrr_ea,

		// 60
		(*Z80).LDrr_hb,
		(*Z80).LDrr_hc,
		(*Z80).LDrr_hd,
		(*Z80).LDrr_he,
		(*Z80).LDrr_hh,
		(*Z80).LDrr_hl,
		(*Z80).LDrHLm_h,
		(*Z80).LDrr_ha,
		(*Z80).LDrr_lb,
		(*Z80).LDrr_lc,
		(*Z80).LDrr_ld,
		(*Z80).LDrr_le,
		(*Z80).LDrr_lh,
		(*Z80).LDrr_ll,
		(*Z80).LDrHLm_l,
		(*Z80).LDrr_la,

		// 70
		(*Z80).LDHLmr_b,
		(*Z80).LDHLmr_c,
		(*Z80).LDHLmr_d,
		(*Z80).LDHLmr_e,
		(*Z80).LDHLmr_h,
		(*Z80).LDHLmr_l,
		(*Z80).HALT,
		(*Z80).LDHLmr_a,
		(*Z80).LDrr_ab,
		(*Z80).LDrr_ac,
		(*Z80).LDrr_ad,
		(*Z80).LDrr_ae,
		(*Z80).LDrr_ah,
		(*Z80).LDrr_al,
		(*Z80).LDrHLm_a,
		(*Z80).LDrr_aa,

		// 80
		(*Z80).ADDr_b,
		(*Z80).ADDr_c,
		(*Z80).ADDr_d,
		(*Z80).ADDr_e,
		(*Z80).ADDr_h,
		(*Z80).ADDr_l,
		(*Z80).ADDHL,
		(*Z80).ADDr_a,
		(*Z80).ADCr_b,
		(*Z80).ADCr_c,
		(*Z80).ADCr_d,
		(*Z80).ADCr_e,
		(*Z80).ADCr_h,
		(*Z80).ADCr_l,
		(*Z80).ADCHL,
		(*Z80).ADCr_a,

		// 90
		(*Z80).SUBr_b,
		(*Z80).SUBr_c,
		(*Z80).SUBr_d,
		(*Z80).SUBr_e,
		(*Z80).SUBr_h,
		(*Z80).SUBr_l,
		(*Z80).SUBHL,
		(*Z80).SUBr_a,
		(*Z80).SBCr_b,
		(*Z80).SBCr_c,
		(*Z80).SBCr_d,
		(*Z80).SBCr_e,
		(*Z80).SBCr_h,
		(*Z80).SBCr_l,
		(*Z80).SBCHL,
		(*Z80).SBCr_a,

		// A0
		(*Z80).ANDr_b,
		(*Z80).ANDr_c,
		(*Z80).ANDr_d,
		(*Z80).ANDr_e,
		(*Z80).ANDr_h,
		(*Z80).ANDr_l,
		(*Z80).ANDHL,
		(*Z80).ANDr_a,
		(*Z80).XORr_b,
		(*Z80).XORr_c,
		(*Z80).XORr_d,
		(*Z80).XORr_e,
		(*Z80).XORr_h,
		(*Z80).XORr_l,
		(*Z80).XORHL,
		(*Z80).XORr_a,

		// B0
		(*Z80).ORr_b,
		(*Z80).ORr_c,
		(*Z80).ORr_d,
		(*Z80).ORr_e,
		(*Z80).ORr_h,
		(*Z80).ORr_l,
		(*Z80).ORHL,
		(*Z80).ORr_a,
		(*Z80).CPr_b,
		(*Z80).CPr_c,
		(*Z80).CPr_d,
		(*Z80).CPr_e,
		(*Z80).CPr_h,
		(*Z80).CPr_l,
		(*Z80).CPHL,
		(*Z80).CPr_a,

		// C0
		(*Z80).RETNZ,
		(*Z80).POPBC,
		(*Z80).JPNZnn,
		(*Z80).JPnn,
		(*Z80).CALLNZnn,
		(*Z80).PUSHBC,
		(*Z80).ADDn,
		(*Z80).RST00,
		(*Z80).RETZ,
		(*Z80).RET,
		(*Z80).JPZnn,
		(*Z80).MAPcb,
		(*Z80).CALLZnn,
		(*Z80).CALLnn,
		(*Z80).ADCn,
		(*Z80).RST08,

		// D0
		(*Z80).RETNC,
		(*Z80).POPDE,
		(*Z80).JPNCnn,
		(*Z80).XX,
		(*Z80).CALLNCnn,
		(*Z80).PUSHDE,
		(*Z80).SUBn,
		(*Z80).RST10,
		(*Z80).RETC,
		(*Z80).RETI,
		(*Z80).JPCnn,
		(*Z80).XX,
		(*Z80).CALLCnn,
		(*Z80).XX,
		(*Z80).SBCn,
		(*Z80).RST18,

		// E0
		(*Z80).LDIOnA,
		(*Z80).POPHL,
		(*Z80).LDIOCA,
		(*Z80).XX,
		(*Z80).XX,
		(*Z80).PUSHHL,
		(*Z80).ANDn,
		(*Z80).RST20,
		(*Z80).ADDSPn,
		(*Z80).JPHL,
		(*Z80).LDmmA,
		(*Z80).XX,
		(*Z80).XX,
		(*Z80).XX,
		(*Z80).ORn,
		(*Z80).RST28,

		// F0
		(*Z80).LDAIOn,
		(*Z80).POPAF,
		(*Z80).LDAIOC,
		(*Z80).DI,
		(*Z80).XX,
		(*Z80).PUSHAF,
		(*Z80).XORn,
		(*Z80).RST30,
		(*Z80).LDHLSPn,
		(*Z80).XX,
		(*Z80).LDAmm,
		(*Z80).EI,
		(*Z80).XX,
		(*Z80).XX,
		(*Z80).CPn,
		(*Z80).RST38}
}
func cbMap() []func(*Z80) {
	return []func(*Z80){
		// CB00
		(*Z80).RLCr_b, (*Z80).RLCr_c, (*Z80).RLCr_d, (*Z80).RLCr_e,
		(*Z80).RLCr_h, (*Z80).RLCr_l, (*Z80).RLCHL, (*Z80).RLCr_a,
		(*Z80).RRCr_b, (*Z80).RRCr_c, (*Z80).RRCr_d, (*Z80).RRCr_e,
		(*Z80).RRCr_h, (*Z80).RRCr_l, (*Z80).RRCHL, (*Z80).RRCr_a,
		// CB10
		(*Z80).RLr_b, (*Z80).RLr_c, (*Z80).RLr_d, (*Z80).RLr_e,
		(*Z80).RLr_h, (*Z80).RLr_l, (*Z80).RLHL, (*Z80).RLr_a,
		(*Z80).RRr_b, (*Z80).RRr_c, (*Z80).RRr_d, (*Z80).RRr_e,
		(*Z80).RRr_h, (*Z80).RRr_l, (*Z80).RRHL, (*Z80).RRr_a,
		// CB20
		(*Z80).SLAr_b, (*Z80).SLAr_c, (*Z80).SLAr_d, (*Z80).SLAr_e,
		(*Z80).SLAr_h, (*Z80).SLAr_l, (*Z80).XX, (*Z80).SLAr_a,
		(*Z80).SRAr_b, (*Z80).SRAr_c, (*Z80).SRAr_d, (*Z80).SRAr_e,
		(*Z80).SRAr_h, (*Z80).SRAr_l, (*Z80).XX, (*Z80).SRAr_a,
		// CB30
		(*Z80).SWAPr_b, (*Z80).SWAPr_c, (*Z80).SWAPr_d, (*Z80).SWAPr_e,
		(*Z80).SWAPr_h, (*Z80).SWAPr_l, (*Z80).XX, (*Z80).SWAPr_a,
		(*Z80).SRLr_b, (*Z80).SRLr_c, (*Z80).SRLr_d, (*Z80).SRLr_e,
		(*Z80).SRLr_h, (*Z80).SRLr_l, (*Z80).XX, (*Z80).SRLr_a,
		// CB40
		(*Z80).BIT0b, (*Z80).BIT0c, (*Z80).BIT0d, (*Z80).BIT0e,
		(*Z80).BIT0h, (*Z80).BIT0l, (*Z80).BIT0m, (*Z80).BIT0a,
		(*Z80).BIT1b, (*Z80).BIT1c, (*Z80).BIT1d, (*Z80).BIT1e,
		(*Z80).BIT1h, (*Z80).BIT1l, (*Z80).BIT1m, (*Z80).BIT1a,
		// CB50
		(*Z80).BIT2b, (*Z80).BIT2c, (*Z80).BIT2d, (*Z80).BIT2e,
		(*Z80).BIT2h, (*Z80).BIT2l, (*Z80).BIT2m, (*Z80).BIT2a,
		(*Z80).BIT3b, (*Z80).BIT3c, (*Z80).BIT3d, (*Z80).BIT3e,
		(*Z80).BIT3h, (*Z80).BIT3l, (*Z80).BIT3m, (*Z80).BIT3a,
		// CB60
		(*Z80).BIT4b, (*Z80).BIT4c, (*Z80).BIT4d, (*Z80).BIT4e,
		(*Z80).BIT4h, (*Z80).BIT4l, (*Z80).BIT4m, (*Z80).BIT4a,
		(*Z80).BIT5b, (*Z80).BIT5c, (*Z80).BIT5d, (*Z80).BIT5e,
		(*Z80).BIT5h, (*Z80).BIT5l, (*Z80).BIT5m, (*Z80).BIT5a,
		// CB70
		(*Z80).BIT6b, (*Z80).BIT6c, (*Z80).BIT6d, (*Z80).BIT6e,
		(*Z80).BIT6h, (*Z80).BIT6l, (*Z80).BIT6m, (*Z80).BIT6a,
		(*Z80).BIT7b, (*Z80).BIT7c, (*Z80).BIT7d, (*Z80).BIT7e,
		(*Z80).BIT7h, (*Z80).BIT7l, (*Z80).BIT7m, (*Z80).BIT7a,
		// CB80
		(*Z80).RES0b, (*Z80).RES0c, (*Z80).RES0d, (*Z80).RES0e,
		(*Z80).RES0h, (*Z80).RES0l, (*Z80).RES0m, (*Z80).RES0a,
		(*Z80).RES1b, (*Z80).RES1c, (*Z80).RES1d, (*Z80).RES1e,
		(*Z80).RES1h, (*Z80).RES1l, (*Z80).RES1m, (*Z80).RES1a,
		// CB90
		(*Z80).RES2b, (*Z80).RES2c, (*Z80).RES2d, (*Z80).RES2e,
		(*Z80).RES2h, (*Z80).RES2l, (*Z80).RES2m, (*Z80).RES2a,
		(*Z80).RES3b, (*Z80).RES3c, (*Z80).RES3d, (*Z80).RES3e,
		(*Z80).RES3h, (*Z80).RES3l, (*Z80).RES3m, (*Z80).RES3a,
		// CBA0
		(*Z80).RES4b, (*Z80).RES4c, (*Z80).RES4d, (*Z80).RES4e,
		(*Z80).RES4h, (*Z80).RES4l, (*Z80).RES4m, (*Z80).RES4a,
		(*Z80).RES5b, (*Z80).RES5c, (*Z80).RES5d, (*Z80).RES5e,
		(*Z80).RES5h, (*Z80).RES5l, (*Z80).RES5m, (*Z80).RES5a,
		// CBB0
		(*Z80).RES6b, (*Z80).RES6c, (*Z80).RES6d, (*Z80).RES6e,
		(*Z80).RES6h, (*Z80).RES6l, (*Z80).RES6m, (*Z80).RES6a,
		(*Z80).RES7b, (*Z80).RES7c, (*Z80).RES7d, (*Z80).RES7e,
		(*Z80).RES7h, (*Z80).RES7l, (*Z80).RES7m, (*Z80).RES7a,
		// CBC0
		(*Z80).SET0b, (*Z80).SET0c, (*Z80).SET0d, (*Z80).SET0e,
		(*Z80).SET0h, (*Z80).SET0l, (*Z80).SET0m, (*Z80).SET0a,
		(*Z80).SET1b, (*Z80).SET1c, (*Z80).SET1d, (*Z80).SET1e,
		(*Z80).SET1h, (*Z80).SET1l, (*Z80).SET1m, (*Z80).SET1a,
		// CBD0
		(*Z80).SET2b, (*Z80).SET2c, (*Z80).SET2d, (*Z80).SET2e,
		(*Z80).SET2h, (*Z80).SET2l, (*Z80).SET2m, (*Z80).SET2a,
		(*Z80).SET3b, (*Z80).SET3c, (*Z80).SET3d, (*Z80).SET3e,
		(*Z80).SET3h, (*Z80).SET3l, (*Z80).SET3m, (*Z80).SET3a,
		// CBE0
		(*Z80).SET4b, (*Z80).SET4c, (*Z80).SET4d, (*Z80).SET4e,
		(*Z80).SET4h, (*Z80).SET4l, (*Z80).SET4m, (*Z80).SET4a,
		(*Z80).SET5b, (*Z80).SET5c, (*Z80).SET5d, (*Z80).SET5e,
		(*Z80).SET5h, (*Z80).SET5l, (*Z80).SET5m, (*Z80).SET5a,
		// CBF0
		(*Z80).SET6b, (*Z80).SET6c, (*Z80).SET6d, (*Z80).SET6e,
		(*Z80).SET6h, (*Z80).SET6l, (*Z80).SET6m, (*Z80).SET6a,
		(*Z80).SET7b, (*Z80).SET7c, (*Z80).SET7d, (*Z80).SET7e,
		(*Z80).SET7h, (*Z80).SET7l, (*Z80).SET7m, (*Z80).SET7a}
}

// End
