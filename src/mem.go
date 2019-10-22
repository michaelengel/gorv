package main

import "fmt"

// little endian only at the moment!

var mem [256*1024]uint32

func r32(addr uint32) uint32 {
	// println("Reading address", addr)
	if (addr & 3) != 0 {
		// exception?
		fmt.Printf("*** word read attempted from unaligned address %x!\n", addr)
	}
	return mem[addr >> 2]
}

func w32(addr uint32, data uint32) {
	// fmt.Printf("*** Writing addr %08x: %08x\n", addr, data);
	mem[addr >> 2] = data
}

func r16(addr uint32) uint16 {
	// TODO: correct?
	// println("Reading address", addr)
	if (addr & 1) != 0 {
		// exception?
		fmt.Printf("*** halfword read attempted from odd address %x!\n", addr)
	}
	return (uint16)(r8(addr >> 2)+r8(addr >> 2+1)<<8)
}

func r8(addr uint32) uint8 {
	fmt.Printf("Reading address %x\n", addr)
	var shift uint32 = uint32(8 * (addr & 3))
	return (uint8)((mem[addr >> 2] >> shift) & 0xff)
}

func w8(addr uint32, data uint8) {
	var shift uint32 = uint32(8 * (addr & 3))
	fmt.Printf("Writing addr %08x with %08x shift %d\n", addr, (uint32(data) << shift), shift)
	mem[addr >> 2] = mem[addr >> 2] | (uint32(data) << shift)
}

