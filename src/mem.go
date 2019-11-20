package main

import "fmt"

// little endian only at the moment!

var mem [256*1024]uint32
var offset uint32=0x80000000 // HACK!

func r32(addr uint32) uint32 {
	addr = addr - offset
	// println("*** Reading address", addr)
	if (addr & 3) != 0 {
		// exception?
		fmt.Printf("*** word read attempted from unaligned address %x!\n", addr)
	}
	return mem[addr >> 2]
}

func w32(addr uint32, data uint32) {
	addr = addr - offset
	fmt.Printf("*** Writing w32 addr %08x: %08x\n", addr, data);
	mem[addr >> 2] = data
}

func r16(addr uint32) uint16 {
	if (addr & 1) != 0 {
		// exception?
		fmt.Printf("*** halfword read attempted from odd address %x!\n", addr)
	}
	fmt.Printf("*** r16: %x %x %x\n", r8(addr), r8(addr+1), (uint16)(r8(addr)+r8(addr+1)<<8))
	return (uint16)(uint16(r8(addr))+uint16(r8(addr+1))<<8) // little endian
}

func w16(addr uint32, data uint16) {
	if (addr & 1) != 0 {
		// exception?
		fmt.Printf("*** halfword write attempted to odd address %x!\n", addr)
	}
	w8(addr, uint8(data & 0xff)) // little endian
	w8(addr+1, uint8(data >> 8))
}

func r8(addr uint32) uint8 {
	var shift uint32 

	addr = addr - offset
	shift = uint32(8 * (addr & 3))

	return (uint8)((mem[addr >> 2] >> shift) & 0xff)
}

func w8(addr uint32, data uint8) {
	addr = addr - offset
	var shift uint32 = uint32(8 * (addr & 3))
	if (addr < 256*1024) {
		mem[addr >> 2] = mem[addr >> 2] & ^(0xff << shift)
		mem[addr >> 2] = mem[addr >> 2] | (uint32(data) << shift)
	} else {
		fmt.Printf("*** byte write to %x: %x\n", addr, data)
	}
}

