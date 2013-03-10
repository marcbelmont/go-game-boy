package main

import "sort"

// /// //
// GPU //
// /// //

type OData struct {
	x, y, tile, palette, yflip, xflip, prio, num int
}

type ODatas []OData

func (s ODatas) Len() int      { return len(s) }
func (s ODatas) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type MySort struct{ ODatas }

func (s MySort) Less(i, j int) bool {
	return s.ODatas[i].num < s.ODatas[j].num
}

type GPU struct {
	tilemap                [][][]byte
	vram, oam, reg         []byte
	objdata, objdatasorted []OData
	bg, obj0, obj1         []byte
	scanrow                []byte

	curline    int
	curscan    int
	linemode   int
	modeclocks uint

	yscrl  int
	xscrl  int
	raster int
	ints   int

	lcdon int
	bgon  int
	objon int
	winon int

	objsize int

	bgtilebase  int
	bgmapbase   int
	wintilebase int

	scrn []uint32
}

func (g *GPU) reset() {
	g.vram = make([]byte, 8192)
	g.oam = make([]byte, 160)
	g.bg = make([]byte, 4)
	g.obj0 = make([]byte, 4)
	g.obj1 = make([]byte, 4)
	for i := 0; i < 4; i++ {
		g.bg[i] = 255
		g.obj0[i] = 255
		g.obj1[i] = 255
	}
	g.tilemap = make([][][]byte, 512)
	for i := 0; i < 512; i++ {
		g.tilemap[i] = make([][]byte, 8)
		for j := 0; j < 8; j++ {
			g.tilemap[i][j] = make([]byte, 8)
		}
	}

	g.reg = make([]byte, 65536)

	g.scrn = make([]uint32, HEIGHT*WIDTH)
	for i := 0; i < len(g.scrn); i++ {
		g.scrn[i] = 0xffffffff
	}

	g.curline = 0
	g.curscan = 0
	g.linemode = 2
	g.modeclocks = 0
	g.yscrl = 0
	g.xscrl = 0
	g.raster = 0
	g.ints = 0

	g.lcdon = 0
	g.bgon = 0
	g.objon = 0
	g.winon = 0

	g.objsize = 0
	g.scanrow = make([]byte, 161)

	g.objdata = make([]OData, 40)
	for i := 0; i < 40; i++ {
		g.objdata[i] = OData{y: -16, x: -8, tile: 0, palette: 0, yflip: 0, xflip: 0, prio: 0, num: i}
	}

	// Set to values expected by BIOS, to start
	g.bgtilebase = 0x0000
	g.bgmapbase = 0x1800
	g.wintilebase = 0x1800
}

func (g *GPU) checkline() {
	g.modeclocks += z80.m
	switch g.linemode {
	// In hblank
	case 0:
		if g.modeclocks >= 51 {
			// End of hblank for last scanline; render screen
			if g.curline == 143 {
				g.linemode = 1
				// todo save image
				mmu.mif |= 1
			} else {
				g.linemode = 2
			}
			g.curline++
			g.curscan += 640
			g.modeclocks = 0
		}
		break

	// In vblank
	case 1:
		if g.modeclocks >= 114 {
			g.modeclocks = 0
			g.curline++
			if g.curline > 153 {
				g.curline = 0
				g.curscan = 0
				g.linemode = 2
			}
		}
		break

	// In OAM-read mode
	case 2:
		if g.modeclocks >= 20 {
			g.modeclocks = 0
			g.linemode = 3
		}
		break

	// In VRAM-read mode
	case 3:
		// Render scanline at end of allotted time
		if g.modeclocks >= 43 {
			g.modeclocks = 0
			g.linemode = 0
			if g.lcdon != 0 {
				if g.bgon != 0 {
					var linebase = g.curscan
					var mapbase = g.bgmapbase + ((((g.curline + g.yscrl) & 255) >> 3) << 5)
					var y = (g.curline + g.yscrl) & 7
					var x = g.xscrl & 7
					var t = (g.xscrl >> 3) & 31
					var w = 160

					if g.bgtilebase != 0 {
						var tile = uint16(g.vram[mapbase+t])
						if tile < 128 {
							tile = 256 + tile
						}
						var tilerow = g.tilemap[tile][y]
						for {
							g.scanrow[160-x] = tilerow[x]
							g.scrn[linebase/4] = buildColor(g.bg[tilerow[x]])
							x++
							if x == 8 {
								t = (t + 1) & 31
								x = 0
								tile = uint16(g.vram[mapbase+t])
								if tile < 128 {
									tile = 256 + tile
								}
								tilerow = g.tilemap[tile][y]
							}
							linebase += 4
							w--
							if w == 0 {
								break
							}
						}
					} else {
						var tilerow = g.tilemap[g.vram[mapbase+t]][y]
						for {
							g.scanrow[160-x] = tilerow[x]
							g.scrn[linebase/4] = buildColor(g.bg[tilerow[x]])
							x++
							if x == 8 {
								t = (t + 1) & 31
								x = 0
								tilerow = g.tilemap[g.vram[mapbase+t]][y]
							}
							linebase += 4
							w--
							if w == 0 {
								break
							}
						}
					}
				}
				if g.objon != 0 {
					var cnt = 0
					// hack
					if g.objsize == 0 && len(g.objdatasorted) != 0 {
						var tilerow, pal []byte
						var x int
						var obj OData
						var linebase = g.curscan
						for i := 0; i < 40; i++ {
							obj = g.objdatasorted[i]
							if obj.y <= g.curline && (obj.y+8) > g.curline {
								if obj.yflip != 0 {
									tilerow = g.tilemap[obj.tile][7-(g.curline-obj.y)]
								} else {
									tilerow = g.tilemap[obj.tile][g.curline-obj.y]
								}
								if obj.palette != 0 {
									pal = g.obj1
								} else {
									pal = g.obj0
								}

								linebase = (g.curline*160 + obj.x) * 4
								if obj.xflip != 0 {
									for x = 0; x < 8; x++ {
										if obj.x+x >= 0 && obj.x+x < 160 {
											if tilerow[7-x] != 0 && (obj.prio != 0 || g.scanrow[x] == 0) {
												g.scrn[linebase/4] = buildColor(pal[tilerow[7-x]])
											}
										}
										linebase += 4
									}
								} else {
									for x = 0; x < 8; x++ {
										if obj.x+x >= 0 && obj.x+x < 160 {
											if tilerow[x] != 0 && (obj.prio != 0 || g.scanrow[x] == 0) {
												g.scrn[linebase/4] = buildColor(pal[tilerow[x]])
											}
										}
										linebase += 4
									}
								}
								cnt++
								if cnt > 10 {
									break
								}
							}
						}
					}
				}
			}
		}
		break
	}
}

func buildColor(b byte) uint32{
	return uint32(b)<<8 + uint32(b)<<16 + uint32(b)<<24 + 0xFF
}

func (g *GPU) updatetile(addr uint16, val uint8) {
	var saddr = addr
	if addr&1 != 0 {
		saddr--
		addr--
	}
	var tile = (addr >> 4) & 511
	var y = (addr >> 1) & 7
	var x byte
	for x = 0; x < 8; x++ {
		var sx byte = 1 << (7 - x)
		var a, b byte
		if (g.vram[saddr] & sx) != 0 {
			a = 1
		}
		if (g.vram[saddr+1] & sx) != 0 {
			b = 2
		}
		g.tilemap[tile][y][x] = a | b
	}
}

func (g *GPU) updateoam(addr uint16, val uint8) {
	addr -= 0xFE00
	var obj = addr >> 2
	if obj < 40 {
		switch addr & 3 {
		case 0:
			g.objdata[obj].y = int(val) - 16
			break
		case 1:
			g.objdata[obj].x = int(val) - 8
			break
		case 2:
			if g.objsize != 0 {
				g.objdata[obj].tile = int(val) & 0xFE
			} else {
				g.objdata[obj].tile = int(val)
			}
			break
		case 3:
			g.objdata[obj].palette = 0
			g.objdata[obj].xflip = 0
			g.objdata[obj].yflip = 0
			g.objdata[obj].prio = 0
			if (val & 0x10) != 0 {
				g.objdata[obj].palette = 1
			}
			if (val & 0x20) != 0 {
				g.objdata[obj].xflip = 1
			}
			if (val & 0x40) != 0 {
				g.objdata[obj].yflip = 1
			}
			if (val & 0x80) != 0 {
				g.objdata[obj].prio = 1
			}
			break
		}
	}

	g.objdatasorted = g.objdata
	sort.Sort(MySort{g.objdatasorted})
}

func (g *GPU) rb(addr uint16) byte {
	var gaddr = addr - 0xFF40
	switch gaddr {
	case 0:
		var res byte = 0
		if g.lcdon != 0 {
			res |= 0x80
		}
		if g.bgtilebase == 0x0000 {
			res |= 0x10
		}
		if g.bgmapbase == 0x1C00 {
			res |= 0x08
		}
		if g.objsize != 0 {
			res |= 0x04
		}
		if g.objon != 0 {
			res |= 0x02
		}
		if g.bgon != 0 {
			res |= 0x01
		}
		return res

	case 1:
		if g.curline == g.raster {
			return byte(4 | g.linemode)
		} else {
			return byte(g.linemode)
		}

	case 2:
		return byte(g.yscrl)

	case 3:
		return byte(g.xscrl)

	case 4:
		return byte(g.curline)

	case 5:
		return byte(g.raster)

	default:
		return byte(g.reg[gaddr])
	}
	return byte(g.reg[gaddr])
}

func (g *GPU) wb(addr uint16, val byte) {
	var gaddr = addr - 0xFF40
	g.reg[gaddr] = val
	switch gaddr {
	case 0:
		if (val & 0x80) != 0 {
			g.lcdon = 1
		} else {
			g.lcdon = 0
		}
		if (val & 0x10) != 0 {
			g.bgtilebase = 0
		} else {
			g.bgtilebase = 0x0800
		}
		if (val & 0x08) != 0 {
			g.bgmapbase = 0x1C00
		} else {
			g.bgmapbase = 0x1800
		}
		if (val & 0x04) != 0 {
			g.objsize = 1
		} else {
			g.objsize = 0
		}
		if (val & 0x02) != 0 {
			g.objon = 1
		} else {
			g.objon = 0
		}
		if (val & 0x01) != 0 {
			g.bgon = 1
		} else {
			g.bgon = 0
		}
		break

	case 2:
		g.yscrl = int(val)
		break

	case 3:
		g.xscrl = int(val)
		break

	case 5:
		g.raster = int(val)

	// OAM DMA
	case 6:
		var v byte
		var i uint16
		for i = 0; i < 160; i++ {
			v = byte(mmu.rb((uint16(val) << 8) + uint16(i)))
			g.oam[i] = v
			g.updateoam(0xFE00+i, v)
		}
		break

	// BG palette mapping
	case 7:
		var i uint16
		for i = 0; i < 4; i++ {
			switch (val >> (i * 2)) & 3 {
			case 0:
				g.bg[i] = 255
			case 1:
				g.bg[i] = 192
			case 2:
				g.bg[i] = 96
			case 3:
				g.bg[i] = 0
			}
		}
		break

	// OBJ0 palette mapping
	case 8:
		var i uint16
		for i = 0; i < 4; i++ {
			switch (val >> (i * 2)) & 3 {
			case 0:
				g.obj0[i] = 255
			case 1:
				g.obj0[i] = 192
			case 2:
				g.obj0[i] = 96
			case 3:
				g.obj0[i] = 0
			}
		}
		break

	// OBJ1 palette mapping
	case 9:
		var i uint16
		for i = 0; i < 4; i++ {
			switch (val >> (i * 2)) & 3 {
			case 0:
				g.obj1[i] = 255
			case 1:
				g.obj1[i] = 192
			case 2:
				g.obj1[i] = 96
			case 3:
				g.obj1[i] = 0
			}
		}
		break
	}
}
