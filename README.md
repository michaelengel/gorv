# gorv
Educational RISC-V emulator written in Go (1.13)

The emulator will eventually emulate the RV32I instruction set.
It also includes a primitive disassembler.

This is a work in progress, only a few instructions are emulated so far.
Expect bugs and horrible code (this is my first Go program, so it will
probably look very much like C ;-)).

Changes:

- 2019-10-18 Implement ELF loading

  Compile with 
  $ riscv32-unknown-elf-gcc -mabi=ilp32 -march=rv32i -nostdlib -e main -o t t.c 
  and pass name of the generated ELF as command line parameter to the emulator

