package main

func main() {
	println("Go RISC-V Emulator")

	loadelf("../test/t1.bin")

	run(0x1000)

}
