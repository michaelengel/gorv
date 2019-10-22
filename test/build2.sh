#!/bin/sh
for i in t2 t3 t4; do
    riscv32-unknown-elf-gcc -mabi=ilp32 -march=rv32im -o $i $i.c -nostdlib -e main
done
