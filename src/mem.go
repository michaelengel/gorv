package main

import "fmt"

// little endian only at the moment!

var mem [256*1024]uint32

func r32(addr uint32) uint32 {
	// println("Reading address", addr)
	return mem[addr]	
}

func w32(addr uint32, data uint32) {
	fmt.Printf("*** Writing addr %08x: %08x\n", addr, data);
	mem[addr] = data
}

func r8(addr uint32) uint8 {
	// println("Reading address", addr)
	var shift uint32 = uint32(8 * (addr & 3))
	return (uint8)((mem[addr] >> shift) & 0xff)
}

func w8(addr uint32, data uint8) {
	var shift uint32 = uint32(8 * (addr & 3))
	addr = addr & 0xfffffffc
	// println("Writing addr", addr, "with", (uint32(data) << shift), "shift", shift)
	mem[addr] = mem[addr] | (uint32(data) << shift)
}

