Go Game Boy
-----------

Based on http://github.com/Two9A/jsGB  

The emulator currently runs opus5.gb and tetris.gb. It has no sound support. Work in progress..

Building the emulator:  
go build src/*.go

Running the emulator:  
./game-boy -cart roms/opus5.gb

Dependencies:  
sdl: go get [github.com/0xe2-0x9a-0x9b/Go-SDL/sdl](http://github.com/0xe2-0x9a-0x9b/Go-SDL/sdl)
