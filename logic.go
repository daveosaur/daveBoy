package main

import (
	"errors"
)

var (
	outOfBounds error = errors.New("program counter out of bounds")
)

func (g *GB) execute(inst byte) (uint16, error) {
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
				addr := getPair(g.H, g.L)
				target = &g.WMem[addr]
			} else {
				target = &g.A
			}
		}
		*target = operand
	//0x8 ADD/ADC
	case 0x8:
		operand := g.getOperand(_low)
		if _low < 0x8 {
			g.addWithMaybeCarry(operand, false)
		} else {
			g.addWithMaybeCarry(operand, true)

		}
	}

	return g.PC, nil
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
		addr := getPair(g.H, g.L)
		operand = g.WMem[addr]
	case 0x7, 0xF:
		operand = g.A
	}
	return operand
}
func (g *GB) addWithMaybeCarry(input byte, carry bool) {
	temp := g.A
	//lazy carry addition
	if carry {
		input++
	}
	g.A += input

	//set flags
	if g.A == 0 {
		g.F.Z = true
	} else {
		g.F.Z = false
	}
	//check halfcarry
	//TODO: verify this i guess?
	if ((g.A&0xF)+(input&0xF))&0x10 == 0x10 {
		g.F.H = true
	} else {
		g.F.H = false
	}
	//check for carry
	//TODO: this probably doesnt actually work. fix later
	if temp > g.A {
		g.F.CY = true
	} else {
		g.F.CY = false
	}
	g.F.N = false //N is always reset
}
