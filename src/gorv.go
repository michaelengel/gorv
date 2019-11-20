package main

import (
        "fmt"
        "os"
)

func main() {
	var elf string = ""
	var sig string = ""

	println("Go RISC-V Emulator")

	argsWithoutProg := os.Args[1:]
	for n, i := range argsWithoutProg {
		if n == 0 {
			fmt.Printf("isa %d %s\n", n, i)
		}
		if n == 1 {
			fmt.Printf("sig %d %s\n", n, i)
			sig = i
			sig = sig[11:]
		}
		if n == 2 {
			elf = i
		}
	}

        if len(os.Args) < 2 {
                fmt.Printf("Usage: %s elf_file\n", os.Args[0])
                os.Exit(1)
        }

	entry, gp, bs, es := loadelf32(elf)

	retval := run(entry, gp)

	fmt.Printf("*** Finished, %d instructions executed\n", ninsns)
	fmt.Printf("main return value: 0x%x = %d\n", retval, retval)

	fmt.Printf("Writing memory %x to %x to signature file %s\n", bs, es, sig)

	f, err := os.Create(sig)
	check(err)

	defer f.Close()

	for i:=bs; i<es; i+=4 {
		s := fmt.Sprintf("%08x\n", r32(i))
		f.WriteString(s)
	}
	f.Sync()
}
