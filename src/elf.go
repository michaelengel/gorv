package main

// TODO: endianness - currently assumes a little endian host machine

import (
	"fmt"
	"io"
	"os"
	"debug/elf"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ioReader(file string) io.ReaderAt {
	r, err := os.Open(file)
	check(err)
	return r
}

func loadelf32(file string) (uint32, uint32, uint32, uint32) {
	
	fmt.Printf("Loading %s\n", file)

	f := ioReader(file)
	_elf, err := elf.NewFile(f)
	check(err)

	// Read and decode ELF identifier
	var ident [16]uint8
	f.ReadAt(ident[0:], 0)
	check(err)
	
	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		fmt.Printf("Bad magic number at %d\n", ident[0:4])
		os.Exit(1)
	}
	
	if _elf.Class.String() != "ELFCLASS32" {
		fmt.Printf("Not a 32 bit ELF file\n")
		os.Exit(1)
	}
	 
	if _elf.Machine.String() != "EM_RISCV" {
		fmt.Printf("Not a RISC V ELF file\n")
		os.Exit(1)
	}
			
	fmt.Printf("Loading ELF sections:\n")
	for n, s := range _elf.Sections {
		fmt.Printf("Section %d Name %s Address %x Offset %x\n", n, s.SectionHeader.Name, s.SectionHeader.Addr, s.SectionHeader.Offset)
	        b1 := make([]byte, s.SectionHeader.Size)

		f.ReadAt(b1, int64(s.SectionHeader.Offset))

		if s.SectionHeader.Addr == 0 {
			continue
		} else if s.SectionHeader.Name == ".sbss" { // clear BSS section
	        	for index, _ := range b1 {
				w8(uint32(s.SectionHeader.Addr + uint64(index)), 0)
			}
		} else {
	        	for index, data := range b1 {
				w8(uint32(s.SectionHeader.Addr + uint64(index)), data)
			}
		}

	}

	var gp uint32 = 0
	var bs uint32 = 0 // begin_signature
	var es uint32 = 0 // end_signature

	fmt.Printf("Loading ELF symbols:\n")
	sy, _ := _elf.Symbols() 
	for index, s := range sy {
		fmt.Printf("Symbol %d name %s\n", index, s.Name)
		if s.Name == "__global_pointer$" {
			fmt.Printf("gp = %x\n", s.Value)
			gp = uint32(s.Value)
		}
		if s.Name == "begin_signature" {
			fmt.Printf("begin_signature = %x\n", s.Value)
			bs = uint32(s.Value)
		}
		if s.Name == "end_signature" {
			fmt.Printf("end_signature = %x\n", s.Value)
			es = uint32(s.Value)
		}
	}

	return uint32(_elf.Entry), uint32(gp), bs, es
}

