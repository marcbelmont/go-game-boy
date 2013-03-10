package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
	"time"
	"flag"
	"fmt"
)

var mmu = MMU{inbios: 1}
var z80 = Z80{}
var gpu = GPU{bgmapbase: 0x1800, wintilebase: 0x1800}
var key = KEY{}
var timer = TIMER{}

func main() {
	initComponents()
	var cartPath = flag.String("cart", "", "cart path")
	flag.Parse()
	if len(*cartPath) != 0 {
		mmu.load(*cartPath)
		runSdl()
	} else {
		flag.PrintDefaults()
		fmt.Println("Keys: arrows, a, r, z, x. Esc to exit")
	}
}

func initComponents() {
	z80.imap = Map()
	z80.cbmap = cbMap()

	mmu.bios = BIOS
	gpu.scrn = make([]uint32, HEIGHT*WIDTH)
	gpu.reset()
	mmu.reset()
	z80.reset()
	key.reset()
	timer.reset()
}

// /// //
// SDL //
// /// //

func runSdl() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}
	var screen = sdl.SetVideoMode(WIDTH, HEIGHT, 32, 0)
	if screen == nil {
		log.Fatal(sdl.GetError())
	}

	ticker := time.NewTicker(time.Second / 50) // 50 Hz

	gbScreen := sdl.CreateRGBSurfaceFrom(
		gpu.scrn,
		WIDTH,
		HEIGHT,
		32,
		WIDTH*4,
		0xff000000,
		0x00ff0000,
		0x0000ff00,
		0x000000ff)

	running := true

	for running {
		select {
		case <-ticker.C:
			for i := 0; i < 10000; i++ {
				z80.exec()
				gpu.checkline()
				timer.inc()
			}
			screen.FillRect(nil, 0xffffff)
			gbScreen = sdl.CreateRGBSurfaceFrom(
				gpu.scrn,
				WIDTH,
				HEIGHT,
				32,
				WIDTH*4,
				0xff000000,
				0x00ff0000,
				0x0000ff00,
				0x000000ff)
			screen.Blit(nil, gbScreen, nil)
			screen.Flip()

		case _event := <-sdl.Events:
			switch e := _event.(type) {
			case sdl.QuitEvent:
				running = false

			case sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				default:
					if e.State == 1 {
						key.keydown(e.Keysym.Sym)
					} else {
						key.keyup(e.Keysym.Sym)
					}
				}

			case sdl.ResizeEvent:
				println("resize screen ", e.W, e.H)

				screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)

				if screen == nil {
					log.Fatal(sdl.GetError())
				}
			}
		}
	}

	gbScreen.Free()
	sdl.Quit()
}
