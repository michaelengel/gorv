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

- 2019-10-18 Clear .sbss section at load time

- 2019-10-22 Fixed memory access bugs, lb/lh sign extension of loaded data
             Don't load ELF sections with offset = 0

- 2019-11-20 gorv now passes 43 of the 55 tests of the riscv-compliance test
             suite (lots of bugs fixed) and implements the command line 
             parameters as used by the test suite for the spike simulator. 

             The remaining failed tests are due to several non-implemented 
             features:

             - CSRxxx and ECALL/EBREAK not yet implemented
             - SRAI/SRLI not yet implemented
             - Tests I-MISALIGN_JMP-01 and I-MISALIGN_LDST-01 fail 
               since exception handling/traps are not yet implemented

             The current version is fixed to executing ELF binaries at
             0x8000_0000 in order to work with the compliance test suite.
             Yes, this is a horrible hack :-).

- 2019-12-17 CSR instructions are implemented 
             gorv now passes 51 of the 55 tests of the riscv-compliance test

             Current deficiencies:
             - ECALL/EBREAK not yet implemented
             - Tests I-MISALIGN_JMP-01 and I-MISALIGN_LDST-01 fail 
               since exception handling/traps are not yet implemented

