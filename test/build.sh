#!/bin/sh
riscv32-unknown-elf-gcc -mabi=ilp32 -march=rv32i -o $1 $1.c -nostdlib -e main
