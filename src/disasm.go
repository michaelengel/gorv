// emulator and disassembler for RV32I instructions

package main

import "fmt"
import "strconv"

// disassemble given instruction
func disasm(insn uint32) (string, byte) {
	var mnem string

	opcode := insn & 0b01111111

	switch opcode { // p16 of RISC V book
		case 0b0110111: 
			mnem = "lui"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b11111111_11111111_11110000_00000000) >> 12
			dis  := join(mnem, reg(rd), strconv.Itoa(to2(off, 20)))
			return dis, 'U'
		case 0b0010111:
			mnem = "auipc"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b11111111_11111111_11110000_00000000) >> 12
			dis  := join(mnem, reg(rd), strconv.Itoa(to2(off, 20)))
			return dis, 'U'
		case 0b1101111: 
			mnem = "jal"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b10000000_00000000_00000000_00000000) >> 31 << 20 | 
			        (insn & 0b01111111_11100000_00000000_00000000) >> 21 <<  1 |
			        (insn & 0b00000000_00010000_00000000_00000000) >> 20 << 11 |
			        (insn & 0b00000000_00001111_11110000_00000000) >> 12 << 12 
			return join(mnem, reg(rd), strconv.Itoa(int(off))), 'J'
		case 0b1100111: 
			mnem = "jalr"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			return join(mnem, reg(rd), reg(rs1), strconv.Itoa(int(off))), 'I'
		case 0b1100011: 
			cond := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  := (insn & 0b00000001_11110000_00000000_00000000) >> 20
			off  := (insn & 0b10000000_00000000_00000000_00000000) >> 31 << 12 | 
			        (insn & 0b01111110_00000000_00000000_00000000) >> 25 <<  5 |
			        (insn & 0b00000000_00000000_00001111_00000000) >>  8 <<  1 |
			        (insn & 0b00000000_00000000_00000000_10000000) >>  7 << 11 
			switch cond {
				case 0b000: mnem = "beq"
				case 0b001: mnem = "bne"
				case 0b100: mnem = "blt"
				case 0b101: mnem = "bge"
				case 0b110: mnem = "bltu"
				case 0b111: mnem = "bgtu"
				default:    mnem = "b??"
			}
			return join(mnem, reg(rs1), reg(rs2), strconv.Itoa(to2(off, 12))), 'B'
		case 0b0000011:
			size := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20

			switch size {
				case 0b000: mnem = "lb"
				case 0b001: mnem = "lh"
				case 0b010: mnem = "lw"
				case 0b100: mnem = "lbu"
				case 0b101: mnem = "lhu"
				default:    mnem = "l??"
			}
                        dis  := join(mnem, reg(rd), strconv.Itoa(to2(off, 11)), "("+reg(rs1)+")")
                        return dis, 'I'
		case 0b0100011:
			size := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  := (insn & 0b00000001_11110000_00000000_00000000) >> 20
			off  := (insn & 0b11111110_00000000_00000000_00000000) >> 25 << 5 | 
			        (insn & 0b00000000_00000000_00001111_10000000) >>  7
          
			switch size {
				case 0b000: mnem = "sb"
				case 0b001: mnem = "sh"
				case 0b010: mnem = "sw"
				default:    mnem = "l??"
			}
                        dis  := join(mnem, reg(rs2), strconv.Itoa(to2(off, 11)), "("+reg(rs1)+")")
                        return dis, 'S'
		case 0b0010011: 
			var off2 int

			aluop :=(insn & 0b00000000_00000000_01110000_00000000) >> 12
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			if (off >= 0b10000000_0000) { 
				off  = off - (1<<11)
				off  = ^off & 0b011111111111
                                off2 = -int(off+1)
			} else {
				off2 = int(off)
			}
			switch aluop {
				case 0b000: mnem = "addi"
				case 0b010: mnem = "slti"
				case 0b011: mnem = "sltiu"
				case 0b100: mnem = "xori"
				case 0b110: mnem = "ori"
				case 0b111: mnem = "andi"
				case 0b001: mnem = "slli"
				case 0b101: mnem = "srli" // or srai...
				default:    mnem = "alu??i"
			}
                        dis  := join(mnem, reg(rd), reg(rs1), strconv.Itoa(off2))
                        return dis, 'I'
		case 0b0110011: 
			aluop  := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			aluop2 := (insn & 0b11111110_00000000_00000000_00000000) >> 25
			rd   :=   (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  :=   (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  :=   (insn & 0b00000001_11110000_00000000_00000000) >> 20
			switch aluop {
				case 0b000: 
					mnem = "alu??"
					if (aluop2 == 0b0000000) { mnem = "add" }
					if (aluop2 == 0b0100000) { mnem = "sub" }
				case 0b001: mnem = "sll"
				case 0b010: mnem = "slt"
				case 0b011: mnem = "sltu"
				case 0b100: mnem = "xor"
				case 0b101:
					mnem = "alu??"
					if (aluop2 == 0b0000000) { mnem = "srl" }
					if (aluop2 == 0b0100000) { mnem = "sra" }
				case 0b110: mnem = "or"
				case 0b111: mnem = "and"
				default:    mnem = "alu??"
			}
                        dis  := join(mnem, reg(rd), reg(rs1), reg(rs2))
                        return dis, 'R'
		case 0b1110011: return "sys", 'I'
		default:        return "unknown", '?'
	}
}

func printregs() {
	fmt.Printf("PC = %08x\n", pc)
	for n, s := range regs {
		fmt.Printf("%s = %08x", reg(uint32(n)), s)
		if (n%8 == 7) {
			fmt.Printf("\n")
		} else {
			fmt.Printf(" ")
		}
	}
}

