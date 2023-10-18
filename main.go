package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	pxSize int32 = 2
	xRes   int32 = 166 * pxSize
	yRes   int32 = 144 * pxSize
)

// figure out a better way to use registers.
// they need to be able to be combined and mutated together.
type GB struct {
	A, B, C, D, E, H, L byte         //general registers
	F                   flagRegister //carry flags and stuff
	PC, SP              uint16       //program counter, stack pointer
	WMem                [8192]byte   //workram!
}

type flagRegister struct {
	Z, N, H, CY bool //zero, negative, half-carry, and carry flags
}

type PPU struct {
	VRAM *[8192]byte //vram!
}

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
	g.PC, err = g.execute(inst)
	if err != nil {
		return err
	}
	return nil
}
func (g *GB) Draw() error {
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
