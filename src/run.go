// emulator for RV32I instructions

package main

import "fmt"

// registers
var regs [32]uint32
var pc   uint32

var ninsns uint32 = 0

func unimp() {
	fmt.Printf("*** unimplemented instruction! ***\n")
}

func illegal() {
	fmt.Printf("*** illegal instruction! ***\n")
}

// emulate given instruction
func emulate(insn uint32) {
	ninsns++

	opcode := insn & 0b00000000_00000000_00000000_01111111

	switch opcode { // p16 of RISC V book
		case 0b0110111: // lui
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b11111111_11111111_11110000_00000000) >> 12
			regs[rd] = off << 12
			pc = pc + 4

		case 0b0010111: // auipc
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b11111111_11111111_11110000_00000000) >> 12
                        off = off << 12
			// regs[rd] = uint32(int(pc) + to2(off, 19))
			regs[rd] = uint32(pc + off)
			pc = pc + 4

		case 0b1101111: // jal
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			off  := (insn & 0b10000000_00000000_00000000_00000000) >> 31 << 20 | 
			        (insn & 0b01111111_11100000_00000000_00000000) >> 21 <<  1 |
			        (insn & 0b00000000_00010000_00000000_00000000) >> 20 << 11 |
			        (insn & 0b00000000_00001111_11110000_00000000) >> 12 << 12 
			if (rd != 0) { regs[rd] = pc + 4 }
			pc = uint32(int(pc) + to2(off, 19))

		case 0b1100111: // jalr
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			fmt.Printf("*** JALR offset = %d\n", to2(off, 11))
			oldpc := pc
			pc = uint32(int(regs[rs1]) + to2(off, 11)) 
			pc = pc & 0xfffffffe
			if (rd != 0) { regs[rd] = oldpc + 4 }

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
				case 0b100: if (to2(regs[rs1], 31) <  to2(regs[rs2], 31)) { jump = true } // blt
				case 0b101: if (to2(regs[rs1], 31) >= to2(regs[rs2], 31)) { jump = true } // bge
				case 0b110: if (regs[rs1] <  regs[rs2]) { jump = true } // bltu
				case 0b111: if (regs[rs1] >= regs[rs2]) { jump = true } // bgtu
				default:    fmt.Printf("Unknown branch condition\n")
			}
			if jump == true {
				pc = uint32(int(pc) + to2(off, 11))
			} else {
				pc = pc + 4
			}
			fmt.Printf("****** new pc = %08x off = %08x = %d\n", pc, off, to2(off, 11))

		case 0b0000011:
			size := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20

			switch size {
				case 0b000: // lb
                        		regs[rd] = uint32(r8(uint32(int(regs[rs1]) + to2(off, 11))))
                        		regs[rd] = uint32(to2(regs[rd], 7))
				case 0b001: // lh
                        		regs[rd] = uint32(r16(uint32(int(regs[rs1]) + to2(off, 11))))
                        		regs[rd] = uint32(to2(regs[rd], 15))
				case 0b010: // lw
                        		regs[rd] = r32(uint32(int(regs[rs1]) + to2(off, 11)))
				case 0b100: // lbu
                        		regs[rd] = uint32(r8(uint32(int(regs[rs1]) + to2(off, 11))))
				case 0b101: // lhu
                        		regs[rd] = uint32(r16(uint32(int(regs[rs1]) + to2(off, 11))))
				default:    // illegal
					illegal()
			}
			pc = pc + 4

		case 0b0100011:
			size := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  := (insn & 0b00000001_11110000_00000000_00000000) >> 20
			off  := (insn & 0b11111110_00000000_00000000_00000000) >> 25 << 5 | 
			        (insn & 0b00000000_00000000_00001111_10000000) >>  7
          
			switch size {
				case 0b000: // sb
                        		w8(uint32(int(regs[rs1]) + to2(off, 11)), uint8(regs[rs2] & 0xff))
				case 0b001: // sh
                        		w16(uint32(int(regs[rs1]) + to2(off, 11)), uint16(regs[rs2] & 0xffff))
				case 0b010: // sw
                        		w32(uint32(int(regs[rs1]) + to2(off, 11)), regs[rs2])
				default:    // illegal
					illegal()
			}
			pc = pc + 4

		case 0b0010011: 
			var off2 int

			aluop :=(insn & 0b00000000_00000000_01110000_00000000) >> 12
			rd   := (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  := (insn & 0b00000000_00001111_10000000_00000000) >> 15
			off  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			offs := (insn & 0b00000001_11110000_00000000_00000000) >> 20 // 5 bit shift offset only
			off2 = to2(off, 11)
			switch aluop {
				case 0b000: regs[rd] = uint32(int(regs[rs1]) + off2) // addi
				case 0b010: if (to2(regs[rs1], 31) < off2) { regs[rd] = 1 } else { regs[rd] = 0 } // slti
				case 0b011: fmt.Printf("*** SLTIU %d %d\n", regs[rs1], off)
					if (regs[rs1] < uint32(off2)) { regs[rd] = 1 } else { regs[rd] = 0 } // sltiu XXX
				case 0b100: regs[rd] = uint32(int(regs[rs1]) ^ off2) // xori
				case 0b110: regs[rd] = uint32(int(regs[rs1]) | off2) // ori
				case 0b111: regs[rd] = uint32(int(regs[rs1]) & off2) // andi
				case 0b001: regs[rd] = regs[rs1] << offs // slli
				case 0b101: unimp() // srli or srai...
				default:    illegal() // illegal ALU op
			}
			pc = pc + 4

		case 0b0110011: 
			aluop  := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			aluop2 := (insn & 0b11111110_00000000_00000000_00000000) >> 25
			rd   :=   (insn & 0b00000000_00000000_00001111_10000000) >>  7
			rs1  :=   (insn & 0b00000000_00001111_10000000_00000000) >> 15
			rs2  :=   (insn & 0b00000001_11110000_00000000_00000000) >> 20
			switch aluop {
				case 0b000: // ALU op
					if (aluop2 == 0b0000000) { regs[rd] = regs[rs1] + regs[rs2] } // add
					if (aluop2 == 0b0100000) { regs[rd] = regs[rs1] - regs[rs2] } // sub
					if (aluop2 == 0b0000001) { regs[rd] = regs[rs1] * regs[rs2] } // mul
				case 0b001: // sll
					regs[rd] = regs[rs1] << to2(regs[rs2] & 0x1f, 5)
				case 0b010: // slt
					if (to2(regs[rs1], 31) < to2(regs[rs2], 31)) { regs[rd] = 1 } else { regs[rd] = 0 } // slt
				case 0b011: // sltu
					if (regs[rs1] < regs[rs2]) { regs[rd] = 1 } else { regs[rd] = 0 } // sltu
				case 0b100: // xor
					regs[rd] = regs[rs1] ^ regs[rs2]
				case 0b101: // shift right log/arith
					if (aluop2 == 0b0000000) { regs[rd] = regs[rs1] >> to2(regs[rs2] & 0x1f, 5) } // srl
					if (aluop2 == 0b0100000) { // sra
						if regs[rs1] & 0x8000_0000 == 0 {
							regs[rd] = regs[rs1] >> to2(regs[rs2] & 0x1f, 5)
						} else {
							x := to2(regs[rs1], 31)
							x = x >> to2(regs[rs2] & 0x1f, 5)
							regs[rd] = uint32(x)
						}
					}
				case 0b110: // or
					regs[rd] = regs[rs1] | regs[rs2]
				case 0b111: // and
					regs[rd] = regs[rs1] & regs[rs2]
				default:
					illegal()
			}
			pc = pc + 4

		case 0b0001111: // FCC
			unimp()
			pc = pc + 4

		case 0b1110011: // CCC
			fct3 := (insn & 0b00000000_00000000_01110000_00000000) >> 12
			imm  := (insn & 0b11111111_11110000_00000000_00000000) >> 20
			unimp()
			pc = pc + 4

		default:        
			fmt.Printf("unknown!\n");
			pc = pc + 4
	}
}

func run(addr uint32, gp uint32) uint32 {
	fmt.Printf("Running from %x with gp %x\n", addr, gp)

	pc = addr
	regs[2] = 0x80010000 // initial sp // HACK
	regs[3] = gp // initial sp

	for {
		insn := r32(pc)

		mnemonic, _ := disasm(insn)
		fmt.Printf("*** %08x: %08x %s\n", pc, insn, mnemonic)

		if (insn == 0xc0001073) { 
			fmt.Printf("*** end emulation ***\n")
			break 
		} // return

		emulate(insn)
		printregs()

		regs[0] = 0 // *HACK* r0 is always zero!

		if (pc == 0) { break } 
	}
	return uint32(regs[10])
}
