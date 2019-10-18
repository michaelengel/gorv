package main

import (
        "fmt"
        "os"
)

func main() {
	println("Go RISC-V Emulator")

        if len(os.Args) < 2 {
                fmt.Println("Usage: elftest elf_file")
                os.Exit(1)
        }

	entry := loadelf32(os.Args[1])

	run(entry)

}
