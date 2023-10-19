package main

import (
	"errors"
)

var (
	outOfBounds error = errors.New("program counter out of bounds")
)

// decode/execute the instruction.
func (g *GB) execute(inst byte) error {

	// high/low bits of instruction to switch on
	_low := inst & 0x0F
	_high := inst >> 4

	switch _high {

	// mostly LD instructions
	// 0x76 is HALT
	case 0x4, 0x5, 0x6, 0x7:
		if inst == 0x76 {
			//halt TODO: probably just return an error w/e
		}
		var target *byte
		operand := g.getOperand(_low)
		switch _high {
		case 0x4:
			if _low < 0x8 {
				target = &g.B
			} else {
				target = &g.C
			}
		case 0x5:
			if _low < 0x8 {
				target = &g.D
			} else {
				target = &g.E
			}
		case 0x6:
			if _low < 0x8 {
				target = &g.H
			} else {
				target = &g.L
			}
		case 0x7:
			if _low < 0x8 {
				//HL fetching requires another cycle
				g.Timer++
				addr := getPair(g.H, g.L)
				target = &g.WMem[addr]
			} else {
				target = &g.A
			}
		}
		*target = operand
		g.Timer++
	//0x8 ADD/ADC
	case 0x8:
		operand := g.getOperand(_low)

		//if _low bit < 0x8 = ADD, otherwise ADC
		g.addWithMaybeCarry(operand, (_low < 0x8))
		g.Timer++
	//0x9 SUB/SBC
	case 0x9:
		operand := g.getOperand(_low)

		//if _low bit < 0x8 = SUB, otherwise SBC
		g.subWithMaybeCarry(operand, (_low < 0x8))
		g.Timer++
	//0xA AND/XOR
	case 0xA:

	//0xB OR/CP
	case 0xB:
	}

	return nil
}

// returns byte at address
func (g *GB) fetch(addr uint16) (byte, error) {
	if g.PC < 0 || int(g.PC)+1 > len(g.WMem) {
		return 0, outOfBounds
	}
	return g.WMem[addr], nil
}

// returns 2 bytes at address
// TODO: probably flip the bytes. wee endian
func (g *GB) doubleFetch(addr uint16) (uint16, error) {
	if g.PC < 0 || int(g.PC)+2 > len(g.WMem) {
		return 0, outOfBounds
	}
	return uint16(uint16(g.WMem[addr])<<8 | uint16(g.WMem[addr+1])), nil
}

// smoosh 2 bytes into uint16
// when reading from memory, flip most/least for endian-ness
func getPair(MostSig, LeastSig byte) uint16 {
	return uint16(uint16(MostSig)<<8 | uint16(LeastSig))
}

// get the operand for instructions 4x-Bx
// may increase timer for HL access
func (g *GB) getOperand(op byte) byte {
	var operand byte
	switch op {
	case 0x0, 0x8:
		operand = g.B
	case 0x1, 0x9:
		operand = g.C
	case 0x2, 0xA:
		operand = g.D
	case 0x3, 0xB:
		operand = g.E
	case 0x4, 0xC:
		operand = g.H
	case 0x5, 0xD:
		operand = g.L
	case 0x6, 0xE:
		//HL fetching requires another cycle
		g.Timer++
		addr := getPair(g.H, g.L)
		operand = g.WMem[addr]
	case 0x7, 0xF:
		operand = g.A
	}
	return operand
}

// ADD/ADC functions
// may increase timer with HL access
func (g *GB) addWithMaybeCarry(inp byte, ADC bool) {
	//input is 1 more if ADC
	if ADC {
		inp++
	}
	//half-carry and carry have to be set first
	g.F.CY = (int16(g.A+inp) > 255)
	//set the low bits of input & A, add together and mask out bit 4
	//then AND with bit 4, half carry if true
	g.F.H = (((g.A & 0x0F) + (inp&0x0F)&0x10) == 0x10)

	g.A += inp

	g.F.Z = (g.A == 0) //zero flag
	g.F.N = false      //N is always reset
}

func (g *GB) subWithMaybeCarry(inp byte, SBC bool) {
	if SBC {
		inp++
	}
	g.F.CY = (int16(g.A-inp) < 0) //lazy borrow check
	//half borrow check
	//TODO: not sure if this actually works but whatever.
	g.F.H = ((g.A & 0x0F) < (inp & 0x0F))

	g.A -= inp

	g.F.Z = (g.A == 0) //zero flag
	g.F.N = true       //N is always set
}
