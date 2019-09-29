// emulator and disassembler for RV32I instructions

package main

import "fmt"
import "strconv"

// registers
var regs [32]uint32
var pc   uint32

// convert bits to 2's complement
func to2(off uint32, bits int) int {
	var off2 int
	if (off >= (1 << bits)) { 
		off  = off - (1<<(bits))
		off  = ^off & ((1<<bits)-1)
       		off2 = -int(off+1)
	} else {
		off2 = int(off)
	}
	return off2
}

// return register name for number
func reg(nr uint32) string {
	r := []string { "zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2", 
			  "fp", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
			  "a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
			  "s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6" }
	return r[int(nr)]
}

// join instruction components for disassembly
func join(str ...string) string {
	var joined string

	for n, s := range str {
		switch n {
			case 0: joined = s
			case 1:  joined = joined + " " + s
			default: joined = joined + ", " + s
		}
	}
	return joined
}

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

// emulate given instruction
func emulate(insn uint32) {
	var mnem string

	opcode := insn & 0b01111111

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
			pc = regs[rd] + uint32(int(pc) + to2(off, 20))

		case 0b1101111: 
			// TODO: signeness, *2???
			mnem = "jal"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b10000000_00000000_00000000_00000000) >> 31 << 20 | 
			        (insn & 0b01111111_11100000_00000000_00000000) >> 21 <<  1 |
			        (insn & 0b00000000_00010000_00000000_00000000) >> 20 << 11 |
			        (insn & 0b00000000_00001111_11110000_00000000) >> 12 << 12 
			if (rd != 0) { regs[rd] = pc + 4 }
			pc = pc + uint32(off)

		case 0b1100111: 
			// TODO: signeness, *2???
			mnem = "jalr"
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			if (rd != 0) { regs[rd] = pc + 4 }
			pc = regs[rs1] + uint32(off)

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
			fmt.Printf("****** new pc = %08x off = %08x off = %d\n", pc, off, to2(off, 12))

		case 0b0000011:
			// TODO: transfer sizes, signedness!
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
                        regs[rd] = r32(uint32(int(regs[rs1]) + to2(off, 11)))
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
				case 0b001: mnem = "sh"
				case 0b010: mnem = "sw"
				default:    mnem = "l??"
			}
                        w32(uint32(int(regs[rs1]) + to2(off, 11)), regs[rs2])
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
				case 0b100: mnem = "xori"
				case 0b110: mnem = "ori"
				case 0b111: mnem = "andi"
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
	println("Running from", addr)

	pc = 0x1000
	regs[2] = 0x100 // initial sp
	regs[8] = 0x100 // initial fp *HACK*

	for {
		insn := r32(pc)

		mnemonic, _ := disasm(insn)
		fmt.Printf("*** %08x: %08x %s\n", pc, insn, mnemonic)

		emulate(insn)
		printregs()

		regs[0] = 0 // *HACK* r0 is always zero!

		if (insn == 0x8067) { break } // return
	}
}
