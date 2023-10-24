package main

import (
	"errors"
)

var (
	outOfBounds     error = errors.New("program counter out of bounds")
	executionHalted error = errors.New("execution halted. probably HALT or STOP")
)

// decode/execute the instruction.
func (g *GB) execute(inst byte) error {

	//16-bit instruction. fetch the next byte and switch on that
	if inst == 0xCB {
		g.PC++ //increment PC to fetch next instruction
		inst2, err := g.fetch(g.PC)
		if err != nil {
			return err
		}
		_low := inst2 & 0x0F
		_high := inst2 >> 4

		switch _high {
		case 0x0: //TODO
		case 0x1:
		case 0x2:
		case 0x3:
		case 0x4:
		case 0x5:
		case 0x6:
		case 0x7:
		case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF: //set bit x to 0 / 1
			targetByte := g.getOperand(inst2)
			var targetBit byte
			switch _high {
			case 0x8, 0xC:
				targetBit = 0x1
				if _low > 0x7 {
					targetBit = 0x2
				}
			case 0x9, 0xD:
				targetBit = 0x4
				if _low > 0x7 {
					targetBit = 0x8
				}
			case 0xA, 0xE:
				targetBit = 0x10
				if _low > 0x7 {
					targetBit = 0x20
				}
			case 0xB, 0xF:
				targetBit = 0x40
				if _low > 0x7 {
					targetBit = 0x80
				}
			}
			if _high < 0xC {
				*targetByte = *targetByte &^ targetBit
			} else {
				*targetByte = *targetByte | targetBit
			}
			g.Cycle++
		}
		return nil

	}
	// high/low nibbles of instruction to switch on
	_low := inst & 0x0F
	_high := inst >> 4

	switch _high {

	//first 4 rows, jumble of instructions organized into columns of 4
	case 0x0, 0x1, 0x2, 0x3:
		switch _low {
		case 0x0:
			switch inst {
			case 0x00: //NOP
				g.Cycle++
			case 0x10: //STOP
				return executionHalted
			case 0x20: //TODO
			case 0x30: //TODO
			}
		case 0x1: //LD d16 into target
			operand, err := g.doubleFetch(g.PC + 1)
			if err != nil {
				return err
			}
			var lowTarget, highTarget *byte
			switch _high {
			case 0x0:
				highTarget, lowTarget = &g.B, &g.C
			case 0x1:
				highTarget, lowTarget = &g.D, &g.E
			case 0x2:
				highTarget, lowTarget = &g.H, &g.L
			case 0x3:
				//this is done different because stack pointer is one 16bit address instead of pairs
				g.SP = operand
				g.Cycle++
				return nil
			}
			*highTarget = byte(operand >> 8)
			*lowTarget = byte(operand & 0x00FF)
			g.Cycle++
		case 0x2: //LD (pair), A
			var highTarget byte
			switch _high {
			case 0x0:
				highTarget = g.Mem[getPair(g.B, g.C)]
			case 0x1:
				highTarget = g.Mem[getPair(g.D, g.E)]
			case 0x2: //HL+
				highTarget = g.Mem[getPair(g.H, g.L)]
				incPair(&g.H, &g.L)
			case 0x3: //HL-
				highTarget = g.Mem[getPair(g.H, g.L)]
				decPair(&g.H, &g.L)
			}
			g.Mem[highTarget] = g.A
			g.Cycle += 2
		case 0x3: //inc register pairs
			switch _high {
			case 0x0:
				incPair(&g.B, &g.C)
			case 0x1:
				incPair(&g.D, &g.E)
			case 0x2:
				incPair(&g.H, &g.L)
			case 0x3: //stack pointer. can just ++
				g.SP++
			}
			g.Cycle += 2
		case 0x4, 0x5: //inc/dec regs
			var target *byte
			switch _high {
			case 0x0:
				target = &g.B
			case 0x1:
				target = &g.D
			case 0x2:
				target = &g.H
			case 0x3:
				target = &g.Mem[getPair(g.H, g.L)]
				g.Cycle += 2
			}
			if _low == 0x4 { //inc on 0x4, dec on 0x5
				g.incByte(target)
			} else {
				g.decByte(target)
			}
			g.Cycle++
		case 0x6: //ld reg, n
			var target *byte
			switch _high {
			case 0x0:
				target = &g.B
			case 0x1:
				target = &g.D
			case 0x2:
				target = &g.H
			case 0x3:
				target = &g.Mem[getPair(g.H, g.L)]
				g.Cycle++
			}
			val, err := g.fetch(g.PC + 1)
			if err != nil {
				return err
			}
			*target = val
			g.Cycle += 2
		case 0x7: //TODO
		case 0x8:
		case 0x9:
		case 0xA:
		case 0xB:
		case 0xC, 0xD: //inc/dec regs part 2
			var target *byte
			switch _high {
			case 0x0:
				target = &g.C
			case 0x1:
				target = &g.E
			case 0x2:
				target = &g.L
			case 0x3:
				target = &g.A
			}
			if _low == 0xC { //inc on 0xC, dec on 0xD
				g.incByte(target)
			} else {
				g.decByte(target)
			}
			g.Cycle++
		case 0xE:
		case 0xF:
		}

	// LD instructions
	// 0x76 is HALT
	case 0x4, 0x5, 0x6, 0x7:
		if inst == 0x76 {
			//halt TODO: probably just return an error w/e
			return executionHalted
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
				g.Cycle++
				addr := getPair(g.H, g.L)
				target = &g.Mem[addr]
			} else {
				target = &g.A
			}
		}
		*target = *operand
		g.Cycle++
	//0x8 ADD/ADC
	case 0x8:
		operand := g.getOperand(_low)

		//if _low bit < 0x8 = ADD, otherwise ADC
		g.addWithMaybeCarry(*operand, (_low < 0x8))
		g.Cycle++
	//0x9 SUB/SBC
	case 0x9:
		operand := g.getOperand(_low)

		//if _low bit < 0x8 = SUB, otherwise SBC
		g.subWithMaybeCarry(*operand, (_low < 0x8))
		g.Cycle++
	//0xA AND/XOR
	case 0xA:
		operand := g.getOperand(_low)

		if _low < 0x8 {
			g.andA(*operand)
		} else {
			g.xorA(*operand)
		}
		g.Cycle++

	//0xB OR/CP
	case 0xB:
		operand := g.getOperand(_low)

		if _low < 0x8 {
			g.orA(*operand)
		} else {
			g.cpA(*operand)
		}
		g.Cycle++
	}

	return nil
}

// returns byte at address
func (g *GB) fetch(addr uint16) (byte, error) {
	if g.PC < 0 || int(g.PC)+1 > len(g.Mem) {
		return 0, outOfBounds
	}
	g.Cycle++
	return g.Mem[addr], nil
}

// returns 2 bytes at address
// adds 2 cycles
// returns in flipped order because little-endian
func (g *GB) doubleFetch(addr uint16) (uint16, error) {
	if g.PC < 0 || int(g.PC)+2 > len(g.Mem) {
		return 0, outOfBounds
	}
	g.Cycle += 2
	return uint16(uint16(g.Mem[addr+1])<<8 | uint16(g.Mem[addr])), nil
}

// smoosh 2 bytes into uint16
// when reading from memory, flip most/least for endian-ness
func getPair(MostSig, LeastSig byte) uint16 {
	return uint16(uint16(MostSig)<<8 | uint16(LeastSig))
}

// get the operand for a chunk of instructions
// may increase timer for HL access
func (g *GB) getOperand(op byte) *byte {
	var operand *byte
	switch op {
	case 0x0, 0x8:
		operand = &g.B
	case 0x1, 0x9:
		operand = &g.C
	case 0x2, 0xA:
		operand = &g.D
	case 0x3, 0xB:
		operand = &g.E
	case 0x4, 0xC:
		operand = &g.H
	case 0x5, 0xD:
		operand = &g.L
	case 0x6, 0xE:
		//HL fetching requires another cycle
		g.Cycle++
		addr := getPair(g.H, g.L)
		operand = &g.Mem[addr]
	case 0x7, 0xF:
		operand = &g.A
	}
	return operand
}

// ADD/ADC functions
// may increase timer with HL access
func (g *GB) addWithMaybeCarry(inp byte, ADC bool) {
	//input is 1 more if ADC
	if ADC && g.F.CY {
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
	if SBC && g.F.CY {
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

// boolean AND input byte with register A
func (g *GB) andA(inp byte) {
	g.A = g.A & inp

	//set flags
	g.F.Z = (g.A == 0)
	g.F.N = false
	g.F.H = true
	g.F.CY = false
}

// boolean XOR input byte with register A
func (g *GB) xorA(inp byte) {
	g.A = g.A ^ inp

	//set flags
	g.F.Z = (g.A == 0)
	g.F.N = false
	g.F.H = false
	g.F.CY = false

}

// boolean OR input byte with register A
func (g *GB) orA(inp byte) {
	g.A = g.A | inp

	//set flags
	g.F.Z = (g.A == 0)
	g.F.N = false
	g.F.H = false
	g.F.CY = false

}

// boolean CP input byte with register A
func (g *GB) cpA(inp byte) {
	result := g.A - inp

	// flagsssss
	g.F.Z = (result == 0)
	g.F.N = false

	// carry/half-carry just stolen from SBC which is already questionable
	// TODO: maybe verify them?
	g.F.CY = (int16(g.A-inp) < 0)
	g.F.H = ((g.A & 0x0F) < (inp & 0x0F))

}

func incPair(high, low *byte) {
	n := getPair(*high, *low)
	n++
	*high = byte(n >> 8)
	*low = byte(n & 0x0F)
}
func decPair(high, low *byte) {
	n := getPair(*high, *low)
	n--
	*high = byte(n >> 8)
	*low = byte(n & 0x0F)
}

// increment byte. sets flags. wee
func (g *GB) incByte(target *byte) {
	*target++

	g.F.Z = (*target == 0)
	g.F.N = false
	//TODO: set half-carry

}

// decrement byte. sets flags. wee
func (g *GB) decByte(target *byte) {
	*target--

	g.F.Z = (*target == 0)
	g.F.N = true
	//TODO: set half-carry
}
