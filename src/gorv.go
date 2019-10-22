package main

import (
        "fmt"
        "os"
)

func main() {
	println("Go RISC-V Emulator")

        if len(os.Args) < 2 {
                fmt.Printf("Usage: %s elf_file\n", os.Args[0])
                os.Exit(1)
        }

	entry := loadelf32(os.Args[1])

	retval := run(entry)

	fmt.Printf("*** Finished, %d instructions executed\n", ninsns)
	fmt.Printf("main return value: 0x%x = %d\n", retval, retval)
}
