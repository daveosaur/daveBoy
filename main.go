package main

//TODO: everythingggggggggggggggg

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	pxSize int32 = 4
	xRes   int32 = 166 * pxSize
	yRes   int32 = 144 * pxSize
)

// gameboy state.
type GB struct {
	A, B, C, D, E, H, L byte         //general registers
	F                   flagRegister //carry flags and stuff
	PC, SP              uint16       //program counter, stack pointer
	WMem                [8192]byte   //workram!
	//gameboy is 4mhz and every operation is a multiple of 4
	//simplify in units of 1
	Cycle uint32 //cycle timer
}

// this is easier than doing bitmasking bleh
type flagRegister struct {
	Z, N, H, CY bool //zero, negative, half-carry, and carry flags
}

type PPU struct {
	VRAM *[8192]byte //vram!
}

// initialize everything
func newGB() (*GB, error) {
	return &GB{
		SP: 0xFFFE,
	}, nil
}

func (g *GB) Update() error {
	inst, err := g.fetch(g.PC)
	if err != nil {
		return err
	}
	err = g.execute(inst)
	if err != nil {
		return err
	}
	return nil
}

// implement everything eventually
func (g *GB) Draw() error {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Gray)

	rl.EndDrawing()
	return nil
}

func main() {
	rl.SetTargetFPS(60)
	rl.InitWindow(xRes, yRes, "daveBoy")

	gb, err := newGB()
	if err != nil {
		panic(err)
	}

	for !rl.WindowShouldClose() {
		gb.Update()
		gb.Draw()

	}
}
