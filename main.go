package main

func main() {
	cpu := NewCPU()
	//cpu.Run()

	cpu.LoadProgram("fixtures/mov_cx_bx.bin")

	cpu.DecodeInatruction()
	cpu.PrintRegisters()
}
