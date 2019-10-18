// utils for emulator and disassembler for RV32I instructions

package main

// convert number in 2's complement format to signed int representation
func to2(off uint32, bits int) int {
	var off2 int
	// fmt.Printf("off: %x bits: %d %x\n", off, bits, (1<<bits))
	// if (off >= (1 << bits)) { 
	if ((off & (1 << (bits))) != 0) { 
		// off  = off - (1<<(bits))
		// off  = ^off & ((1<<bits)-1)
		off  = off ^ ((1<<(bits+1))-1) // invert all bits of number incl. sign
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

