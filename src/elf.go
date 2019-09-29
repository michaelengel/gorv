package main

import (
    "os"
)

// just fake elf at the moment, read objcopy'd binaries

func loadelf(file string) {
	var addr uint32

	// reset vector is 0x0, start code at 0x1000
	w32(0x0, 0x1000)

	f, _ := os.Open(file)

	addr = 0x1000

	b1 := make([]byte, 1000)
	f.Read(b1)

	for index, data := range b1 {
		w8(addr + uint32(index), data)
	}
}
