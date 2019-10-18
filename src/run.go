// emulator and disassembler for RV32I instructions

package main

import "fmt"
import "strconv"

// registers
var regs [32]uint32
var pc   uint32

// emulate given instruction
func emulate(insn uint32) {
	var mnem string

	opcode := insn & 0b00000000_00000000_00000000_01111111

	switch opcode { // p16 of RISC V book
		case 0b0110111: 
			mnem = "lui"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b11111111_11111111_11110000_00000000) >> 12
			regs[rd] = off << 12
			pc = pc + 4

		case 0b0010111:
			mnem = "auipc"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b11111111_11111111_11110000_00000000) >> 12
                        off = off << 12
			regs[rd] = uint32(int(pc) + to2(off, 20))

		case 0b1101111: 
			mnem = "jal"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b10000000_00000000_00000000_00000000) >> 31 << 20 | 
			        (insn & 0b01111111_11100000_00000000_00000000) >> 21 <<  1 |
			        (insn & 0b00000000_00010000_00000000_00000000) >> 20 << 11 |
			        (insn & 0b00000000_00001111_11110000_00000000) >> 12 << 12 
			if (rd != 0) { regs[rd] = pc + 4 }
			pc = uint32(int(pc) + to2(off, 20))

		case 0b1100111: 
			mnem = "jalr"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			if (rd != 0) { regs[rd] = pc + 4 }
			pc = uint32(int(regs[rs1]) + to2(off, 20))

		case 0b1100011: 
			var jump bool = false
			cond := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  := (insn & 0b00000001_11110000_00000000_00000000) >> 20
			off  := (insn & 0b10000000_00000000_00000000_00000000) >> 31 << 12 | 
			        (insn & 0b01111110_00000000_00000000_00000000) >> 25 <<  5 |
			        (insn & 0b00000000_00000000_00001111_00000000) >>  8 <<  1 |
			        (insn & 0b00000000_00000000_00000000_10000000) >>  7 << 11 
			switch cond {
				// TODO: 2's complement!
				case 0b000: if (regs[rs1] == regs[rs2]) { jump = true } // beq
				case 0b001: if (regs[rs1] != regs[rs2]) { jump = true } // bne
				case 0b100: if (regs[rs1] <  regs[rs2]) { jump = true } // blt
				case 0b101: if (regs[rs1] >= regs[rs2]) { jump = true } // bge
				case 0b110: if (regs[rs1] <  regs[rs2]) { jump = true } // bltu
				case 0b111: if (regs[rs1] >= regs[rs2]) { jump = true } // bgtu
				default:    fmt.Printf("Unknown branch condition\n")
			}
			if jump == true {
				pc = uint32(int(pc) + to2(off, 12))
			} else {
				pc = pc + 4
			}
			fmt.Printf("****** new pc = %08x off = %08x = %d\n", pc, off, to2(off, 12))

		case 0b0000011:
			// TODO: transfer sizes, signedness!
			size := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20

			switch size {
				case 0b000: mnem = "lb"
                        			regs[rd] = uint32(r8(uint32(int(regs[rs1]) + to2(off, 12))))
				case 0b001: mnem = "lh"
                        			// regs[rd] = uint32(r16(uint32(int(regs[rs1]) + to2(off, 12))))
				case 0b010: mnem = "lw"
                        			regs[rd] = r32(uint32(int(regs[rs1]) + to2(off, 12)))
				case 0b100: mnem = "lbu"
				case 0b101: mnem = "lhu"
				default:    mnem = "l??"
			}
			pc = pc + 4

		case 0b0100011:
			// TODO: transfer sizes, signedness!
			size := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  := (insn & 0b00000001_11110000_00000000_00000000) >> 20
			off  := (insn & 0b11111110_00000000_00000000_00000000) >> 25 << 5 | 
			        (insn & 0b00000000_00000000_00001111_10000000) >>  7
          
			switch size {
				case 0b000: mnem = "sb"
                        			w8(uint32(int(regs[rs1]) + to2(off, 12)), uint8(regs[rs2] & 0xff))
				case 0b001: mnem = "sh"
                        			// w16(uint32(int(regs[rs1]) + to2(off, 12)), uint16(regs[rs2] & 0xffff))
				case 0b010: mnem = "sw"
                        			w32(uint32(int(regs[rs1]) + to2(off, 12)), regs[rs2])
				default:    mnem = "l??"
			}
			pc = pc + 4

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
				case 0b000: regs[rd] = uint32(int(regs[rs1]) + off2) // addi
				case 0b010: mnem = "slti"
				case 0b011: mnem = "sltiu"
				case 0b100: regs[rd] = uint32(int(regs[rs1]) ^ off2) // xori
				case 0b110: regs[rd] = uint32(int(regs[rs1]) | off2) // ori
				case 0b111: regs[rd] = uint32(int(regs[rs1]) & off2) // andi
				case 0b001: mnem = "slli"
				case 0b101: mnem = "srli" // or srai...
				default:    mnem = "alu??i"
			}
                        dis  := join(mnem, reg(rd), reg(rs1), strconv.Itoa(off2))
			fmt.Printf("", dis)
			pc = pc + 4

		case 0b0110011: 
			aluop  := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			aluop2 := (insn & 0b11111110_00000000_00000000_00000000) >> 25
			rd   :=   (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  :=   (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  :=   (insn & 0b00000001_11110000_00000000_00000000) >> 20
			switch aluop {
				case 0b000: 
					mnem = "alu??"
					if (aluop2 == 0b0000000) { regs[rd] = regs[rs1] + regs[rs2] } // add
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
			fmt.Printf("", dis)
			pc = pc + 4

		case 0b1110011: 
			pc = pc + 4

		default:        
			fmt.Printf("unknown!\n");
			pc = pc + 4
	}
}

func run(addr uint32) {
	fmt.Printf("Running from %x\n", addr)

	pc = addr
	regs[2] = 0x10000 // initial sp
	regs[8] = 0x10000 // initial fp *HACK*

	for {
		insn := r32(pc)

		mnemonic, _ := disasm(insn)
		fmt.Printf("*** %08x: %08x %s\n", pc, insn, mnemonic)

		emulate(insn)
		printregs()

		regs[0] = 0 // *HACK* r0 is always zero!

		// if (insn == 0x8067) { break } // return
		if (pc == 0) { break } 
	}
}
