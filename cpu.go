package main

import (
	"fmt"
	"io"
	"os"
)

// i8086 CPU
type CPU struct {
	PC uint16 // Program Counter
	SP uint16 // Stack Pointer
	AX uint16 // Accumulator
	BX uint16 // Base
	CX uint16 // Counter
	DX uint16 // Data
	SI uint16 // Source Index
	DI uint16 // Destination Index
	BP uint16 // Base Pointer
	CS uint16 // Code Segment
	DS uint16 // Data Segment
	ES uint16 // Extra Segment
	SS uint16 // Stack Segment
	IP uint16 // Instruction Pointer
	FL uint16 // Flag Register

	Flag uint16

	programSize int

	// 1MB of memory
	Memory [1048576]byte
}

func getAL(c *CPU) uint8 {
	return uint8(c.AX & 0xFF)
}

func getAH(c *CPU) uint8 {
	return uint8(c.AX >> 8)
}

func getBL(c *CPU) uint8 {
	return uint8(c.BX & 0xFF)
}

func getBH(c *CPU) uint8 {
	return uint8(c.BX >> 8)
}

func getCL(c *CPU) uint8 {
	return uint8(c.CX & 0xFF)
}

func getCH(c *CPU) uint8 {
	return uint8(c.CX >> 8)
}

func getDL(c *CPU) uint8 {
	return uint8(c.DX & 0xFF)
}

func getDH(c *CPU) uint8 {
	return uint8(c.DX >> 8)
}

func setAL(c *CPU, v uint8) {
	c.AX = (c.AX & 0xFF00) | uint16(v)
}

func setAH(c *CPU, v uint8) {
	c.AX = (c.AX & 0x00FF) | (uint16(v) << 8)
}

func setBL(c *CPU, v uint8) {
	c.BX = (c.BX & 0xFF00) | uint16(v)
}

func setBH(c *CPU, v uint8) {
	c.BX = (c.BX & 0x00FF) | (uint16(v) << 8)
}

func setCL(c *CPU, v uint8) {
	c.CX = (c.CX & 0xFF00) | uint16(v)
}

func setCH(c *CPU, v uint8) {
	c.CX = (c.CX & 0x00FF) | (uint16(v) << 8)
}

func setDL(c *CPU, v uint8) {
	c.DX = (c.DX & 0xFF00) | uint16(v)
}

func setDH(c *CPU, v uint8) {
	c.DX = (c.DX & 0x00FF) | (uint16(v) << 8)
}

func (c *CPU) PrintRegisters() {
	fmt.Printf("AX: %04X %016b AH: %08b AL: %08b\n", c.AX, c.AX, getAH(c), getAL(c))
	fmt.Printf("BX: %04X %016b BH: %08b BL: %08b\n", c.BX, c.BX, getBH(c), getBL(c))
	fmt.Printf("CX: %04X %016b CH: %08b CL: %08b\n", c.CX, c.CX, getCH(c), getCL(c))
	fmt.Printf("DX: %04X %016b DH: %08b DL: %08b\n", c.DX, c.DX, getDH(c), getDL(c))
	fmt.Printf("SI: %04X %016b\n", c.SI, c.SI)
	fmt.Printf("DI: %04X %016b\n", c.DI, c.DI)
	fmt.Printf("BP: %04X %016b\n", c.BP, c.BP)
	fmt.Printf("CS: %04X %016b\n", c.CS, c.CS)
	fmt.Printf("DS: %04X %016b\n", c.DS, c.DS)
	fmt.Printf("ES: %04X %016b\n", c.ES, c.ES)
	fmt.Printf("SS: %04X %016b\n", c.SS, c.SS)
	fmt.Printf("IP: %04X %016b\n", c.IP, c.IP)
	fmt.Printf("FL: %04X %016b\n", c.FL, c.FL)
	fmt.Printf("SP: %04X %016b\n", c.SP, c.SP)
	fmt.Printf("PC: %04X %016b\n", c.PC, c.PC)

	fmt.Printf("Flag: %04X %016b\n", c.Flag, c.Flag)

	c.PrintMemory()
}

func (c *CPU) PrintMemory() {
	fmt.Printf("Memory:\n")

	// print binary
	for i := 0; i < c.programSize; i++ {
		fmt.Printf("%08b\n", c.Memory[i])
	}

	fmt.Printf("\n")

	for i := 0; i < c.programSize; i += 16 {

		// Print Hex
		fmt.Printf("%04X: ", i)
		for j := 0; j < 16; j++ {
			fmt.Printf("%02X ", c.Memory[i+j])
		}

		fmt.Printf(" | ")

		// Print ASCII
		for j := 0; j < 16; j++ {
			if c.Memory[i+j] >= 0x20 && c.Memory[i+j] <= 0x7E {
				fmt.Printf("%c", c.Memory[i+j])
				continue
			}
			fmt.Printf(".")
		}

		fmt.Printf("\n")
	}
}

type Instruction struct {
	// mnemonic
	// destination
	// source
	// code   -> 6 bits
	// D      -> 1 bit
	// W      -> 1 bit
	// Mod    -> 2 bits
	// Reg    -> 3 bits
	// R/M    -> 3 bits

}

/*

HARDWARE REFERENCE INFORMATION
8086 family uses manual page 160


Single-Bit Field Encoding

	Field |		Value	  |	Notes
	------|---------------|---------------------------
	   S  |		0		  |	No sign extension
	   S  |		1		  |	Sign extend 8-bit immediate data to 16 bits if W=1
	   W  |		0		  | Instruction operateson bytedata
	   W  |		1		  | Instruction operates on word data
	   D  |		0		  | Instruction source is specified in REG field
	   D  |		1		  | Instruction destination is specified in REG field
	   V  |		0		  | Shift/ rotate count is one
	   V  |		1		  | Shift/rotate count is specified in CL register
	   Z  |		0		  | Repeat/loop whilezero flag is clear
	   Z  |		1		  | Repeat/loop whilezero flag is set



REG (Register) Field Encoding

	REG  |	w = 0 	 |	w = 1
	-----|-----------|-----------
	000  |	AL		 |	AX
	001  |	CL		 |	CX
	010  |	DL		 |	DX
	011  |	BL		 |	BX
	100  |	AH		 |	SP
	101  |	CH		 |	BP
	110  |	DH		 |	SI
	111  |	BH		 |	DI


MOD (Mode) Field Encoding

	MOD  |		NOTES
	-----|---------------------------
	00   |		Memory Mode, no displacement follows*
	01   |		Memory Mode, 8-bit displacement follows
	10   |		Memory Mode, 16-bit displacement follows
	11   |		Register Mode

	*Except when R/M = 110, then 16-bit displacement follows


R/M (Register/Memory) Field Encoding

	MOD = 11		     |		EFFECTIVE ADDRESS CALCULATION
	---------------------|-------------------------------------------------------
	R/M	 | w = 0 | w = 1 |	R/M		|	MOD = 00	|	MOD = 01	|	MOD = 10
	-----|-------|-------|----------|---------------|---------------|------------
	000  |	AL	 |	AX	 |	000	    |	(BX)+(SI)	| (BX)+(SI)+d8	| (BX)+(SI)+d16
	001  |	CL	 |	CX	 |	001	    |	(BX)+(DI)	| (BX)+(DI)+d8	| (BX)+(DI)+d16
	010  |	DL	 |	DX	 |	010	    |	(BP)+(SI)	| (BP)+(SI)+d8	| (BP)+(SI)+d16
	011  |	BL	 |	BX	 |	011	    |	(BP)+(DI)	| (BP)+(DI)+d8	| (BP)+(DI)+d16
	100  |	AH	 |	SP	 |	100	    |	(SI)		| (SI)+d8		| (SI)+d16
	101  |	CH	 |	BP	 |	101	    |	(DI)		| (DI)+d8		| (DI)+d16
	110  |	DH	 |	SI	 |	110	    |	(BP)		| (BP)+d8		| (BP)+d16
	111  |	BH	 |	DI	 |	111	    |	(BX)		| (BX)+d8		| (BX)+d16

*/

var (
	mnemonics = map[uint8]string{
		0b100010: "MOV",
	}
)

func (c *CPU) calcLen(opcode uint8, d uint8, w uint8, mod uint8, rm uint8) (uint8, error) {
	length := uint8(0)
	switch opcode {
	case 0b100010:
		length = 1
		if w == 1 {
			length = 2
		}
	default:
		return 0, fmt.Errorf("invalid opcode: %d", opcode)

	}

	return length, nil
}

func (c *CPU) DecodeInatruction() (Instruction, error) {

	// Fetch
	memory := c.Memory[c.PC]
	c.PC++

	// Decode
	opcode := (memory & 0xFC) >> 2 // 6 bits -> opcode
	d := (memory & 0x2) >> 1       // 1 bit -> destination or source
	w := (memory & 0x1)            // 1 bit -> word or byte

	// Fetch mod reg r/m
	memory = c.Memory[c.PC]
	c.PC++

	// Decode
	mod := (memory & 0xC0) >> 6 // 2 bits -> mode
	reg := (memory & 0x38) >> 3 // 3 bits -> register
	rm := (memory & 0x07)       // 3 bits -> register or memory

	// Print Instruction
	fmt.Printf("menemonic: %s\n", mnemonics[opcode])

	// print binary
	fmt.Printf("opcode: %06b d: %01b w: %01b mod: %02b reg: %03b rm: %03b\n", opcode, d, w, mod, reg, rm)

	return Instruction{}, nil
}

func (c *CPU) LoadProgram(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	i := 0
	for ; i < len(b); i++ {
		c.Memory[i] = b[i]
	}

	c.programSize = i
	return nil
}

func (c *CPU) Run() error {
	for {
		// Fetch
		//opcode := c.Memory[c.PC]

	}

	return nil
}

func NewCPU() *CPU {
	return &CPU{}
}
