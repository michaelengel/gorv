#!/bin/sh
riscv32-unknown-elf-gcc -mabi=ilp32 -march=rv32i -c t1.c 
riscv32-unknown-elf-objcopy -O binary t1.o t1.bin
