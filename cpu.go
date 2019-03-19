package main

import "C"
import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// RunState specifies the current running state of the Processor.
type RunState uint32

const (
	// RunStateStopped indicates that the Processor is no longer incrementing the
	// program counter and executing instructions.
	RunStateStopped RunState = iota

	// RunStateRunning indicates that the Processor is currently executing
	// instructions until interrupted.
	RunStateRunning
)

// CPU is a Processor to emulate the LC-3 CPU.
type CPU struct {
	Reg    [8]uint16     // registers
	PC     uint16        // Program Counter
	Memory [65536]uint16 // CPU Memory
	//Memory       []byte        // CPU Memory
	CondRegister *CondRegister // Condition Flags Register

	TimerStarted bool
	TimerStart   time.Time

	OP       uint16   // current opcode
	runState RunState // current state
}

// CondRegister stores the state of the CPU condition flags register.
type CondRegister struct {
	P bool // Sign (S), set if the result is negative.
	Z bool // Zero (Z), set if the result is zero.
	N bool // Parity (P), set if the number of 1 bits in the result is even.
}

// Memory Mapped Registers
const (
	// Keyboard status
	MemRegKBSR uint16 = 0xFE00

	// Keyboard data
	MemRegKBDR uint16 = 0xFE02
)

// List of OpCodes
const (
	OpBR   uint16 = iota // branch
	OpADD                // add
	OpLD                 // load
	OpST                 // store
	OpJSR                // jump register
	OpAND                // bitwise and
	OpLDR                // load register
	OpSTR                // store register
	OpRTI                // unused
	OpNOT                // bitwise not
	OpLDI                // load indirect
	OpSTI                // store indrect
	OpJMP                // jump
	OpRES                // reserved (unused)
	OpLEA                // load effective address
	OpTRAP               // execute trap
)

// List of Trap codes
const (
	TrapGETC  uint16 = 0x20 // get character from keyboard
	TrapOUT   uint16 = 0x21 // output a character
	TrapPUTS  uint16 = 0x22 // output a word string
	TrapIN    uint16 = 0x23 // input a string
	TrapPUTSP uint16 = 0x24 // output a byte string
	TrapHALT  uint16 = 0x25 // halt the program
)

// NewCPU creates a new instance of the CPU
func NewCPU() *CPU {
	cpu := CPU{}
	return &cpu
}

// Run executes any program loaded into memory, starting from the program
// counter value, running until completion.
func (c *CPU) Run() (err error) {
	//var cycles uint8 = 4

	if len(c.Memory) == 0 {
		return errNoProgram
	}

	fmt.Println("Loaded memory: ", len(c.Memory))

	for {
		//c.ProcessInput()
		err = c.Step()
		if err != nil || c.runState == RunStateStopped {
			break
		}
	}

	return
}

// Reset the CPU
func (c *CPU) Reset() {
	// set the PC to the starting position
	// 0x3000 is the default
	c.PC = 0x3000

	// Reset the condition register flags
	c.CondRegister = &CondRegister{}
}

// Step executes the program loaded into memory
func (c *CPU) Step() (err error) {
	c.runState = RunStateRunning
	//fmt.Println("PC: ", c.PC)
	c.EmulateInstruction()
	//Increment MCC
	c.Memory[0xFFFF]++
	return
}

// Stop instructs the processor to stop processing instructions.
func (c *CPU) Stop() (err error) {
	c.runState = RunStateStopped
	return
}

// ProcessInput handles keyboard input
func (c *CPU) ProcessInput() (err error) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	var b = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		log.Println(fmt.Sprintf("0x%04x: Got byte: %s", c.PC, string(b)))
	}
	return
}

// ReadMemory reads an address from memory
func (c *CPU) ReadMemory(address uint16) uint16 {
	//log.Printf("Reading memory address: 0x%04X", address)
	if address == MemRegKBSR {
		// get ascii key code
		ascii := 0xFF87 // yes
		//ascii, _, _ := getChar()
		//log.Println(fmt.Sprintf("0x%04x: getChar %d", c.PC, ascii))
		if ascii > 0 { // TODO - should be triggered by key down
			c.Memory[MemRegKBSR] = (1 << 15)
			c.Memory[MemRegKBDR] = uint16(ascii)
		} else {
			c.Memory[MemRegKBSR] = 0
		}
	}

	switch {
	case address <= 65535:
		//log.Printf("Value is: %d", c.Memory[address])
		return uint16(c.Memory[address])
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

// WriteMemory writes to an address in memory
func (c *CPU) WriteMemory(address uint16, value uint16) {
	switch {
	case address <= 65535:
		c.Memory[address] = value
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

// EmulateInstruction emulates the LC-3 instruction
func (c *CPU) EmulateInstruction() (err error) {
	var pc uint16 = c.PC + 1

	instr := c.ReadMemory(c.PC)
	op := instr >> 12
	//fmt.Printf("Received Inst:0x%04x Op:%d\n", instr, op)

	// process the current opcode
	switch op {
	case OpBR:
		n := extract1C(instr, 11, 11) == 1
		z := extract1C(instr, 10, 10) == 1
		p := extract1C(instr, 9, 9) == 1
		PCoffset9 := extract2C(instr, 8, 0)

		brString := fmt.Sprintf("0x%04x: BR", c.PC)
		if n {
			brString += fmt.Sprintf("n")
		}
		if z {
			brString += fmt.Sprintf("z")
		}
		if p {
			brString += fmt.Sprintf("p")
		}
		brString += fmt.Sprintf(" #%d\n", int16(PCoffset9))
		//log.Println(brString)

		if (n && c.CondRegister.N) || (z && c.CondRegister.Z) || (p && c.CondRegister.P) {
			pc += PCoffset9
		}
	case OpJMP:
		baseR := extract1C(instr, 8, 6)
		pc = c.Reg[baseR]
	case OpADD:
		dr := extract1C(instr, 11, 9)
		sr1 := extract1C(instr, 8, 6)
		bit5 := extract1C(instr, 5, 5)
		if bit5 == 1 {
			imm5 := extract2C(instr, 4, 0)
			//log.Println("0x%04x: ADD R%d,R%d,#%d\n", c.PC, dr, sr1, int16(imm5))
			c.Reg[dr] = c.Reg[sr1] + imm5
		} else {
			sr2 := extract1C(instr, 2, 0)
			//log.Println("0x%04x: ADD R%d,R%d,#%d\n", c.PC, dr, sr1, sr2)
			c.Reg[dr] = c.Reg[sr1] + c.Reg[sr2]
		}
		c.SetCC(dr)
	case OpAND:
		dr := extract1C(instr, 11, 9)
		sr1 := extract1C(instr, 8, 6)
		bit5 := extract1C(instr, 5, 5)
		if bit5 == 1 {
			imm5 := extract2C(instr, 4, 0)
			c.Reg[dr] = c.Reg[sr1] & imm5
		} else {
			sr2 := extract1C(instr, 2, 0)
			c.Reg[dr] = c.Reg[sr1] & c.Reg[sr2]
		}
		c.SetCC(c.Reg[dr])
	case OpNOT:
		dr := extract1C(instr, 11, 9)
		sr := extract1C(instr, 8, 6)
		c.Reg[dr] = ^c.Reg[sr]
		c.SetCC(c.Reg[dr])
	case OpLD:
		dr := extract1C(instr, 11, 9)
		PCoffset9 := extract2C(instr, 8, 0)
		c.Reg[dr] = c.ReadMemory(c.PC + PCoffset9)
		c.SetCC(c.Reg[dr])
		//log.Println(fmt.Sprintf("0x%04x: LD R%d,%d", c.PC, dr, PCoffset9))
	case OpLDI:
		dr := extract1C(instr, 11, 9)
		PCoffset9 := extract2C(instr, 8, 0)
		addr := c.ReadMemory(pc + PCoffset9)
		c.Reg[dr] = c.ReadMemory(addr)
		c.SetCC(c.Reg[dr])
		//log.Println(fmt.Sprintf("0x%04x: LDI R%d,0x%04x", c.PC, dr, addr))
	case OpJSR:
		bit11 := extract1C(instr, 11, 11)
		c.Reg[7] = pc
		if bit11 == 1 {
			PCoffset11 := extract2C(instr, 10, 0)
			pc += PCoffset11
			//log.Println(fmt.Sprintf("0x%04x: JSR BIT1 0x%04x,%d", c.PC, c.Reg[7], PCoffset11))
		} else {
			baseR := extract2C(instr, 8, 6)
			pc = c.Reg[baseR]
			//log.Println(fmt.Sprintf("0x%04x: JSR BASER 0x%04x,%d", c.PC, c.Reg[7], baseR))
		}
	case OpLDR:
		dr := extract1C(instr, 11, 9)
		baseR := extract1C(instr, 8, 6)
		offset6 := extract2C(instr, 5, 0)
		c.Reg[dr] = c.ReadMemory(c.Reg[baseR] + offset6)
		c.SetCC(c.Reg[dr])
	case OpLEA:
		dr := extract1C(instr, 11, 9)
		PCoffset9 := extract2C(instr, 8, 0)
		c.Reg[dr] = c.PC + PCoffset9
		c.SetCC(c.Reg[dr])
		//log.Println(fmt.Sprintf("0x%04x: LEA R%d,%d", c.PC, dr, PCoffset9))
	case OpST:
		sr := extract1C(instr, 11, 9)
		PCoffset9 := extract2C(instr, 8, 0)
		c.WriteMemory(pc+PCoffset9, c.Reg[sr])
		//log.Println(fmt.Sprintf("0x%04x: ST R%d,%d", c.PC, sr, PCoffset9))
	case OpSTI:
		sr := extract1C(instr, 11, 9)
		PCoffset9 := extract2C(instr, 8, 0)
		c.WriteMemory(c.ReadMemory(pc+PCoffset9), c.Reg[sr])
	case OpSTR:
		sr := extract1C(instr, 11, 9)
		baseR := extract1C(instr, 8, 6)
		offset6 := extract2C(instr, 5, 0)
		c.WriteMemory(c.Reg[baseR]+offset6, c.Reg[sr])
		//log.Println(fmt.Sprintf("0x%04x: STR R%d 0x%04x,0x%04x", c.PC, sr, c.Reg[baseR]+offset6, c.Reg[sr]))
	case OpTRAP:
		trapCode := instr & 0xFF
		switch trapCode {
		case TrapGETC:
			log.Fatalf("GETC Trap code not implemented: 0x%04X", instr)
			// read a single ASCII character
			ascii, _, _ := getChar()
			c.Reg[0] = uint16(ascii)
		case TrapOUT:
			//fmt.Println("trapout")
			chr := rune(c.Reg[0])
			fmt.Printf("%c", chr)
			//fmt.Printf("%c\n", ch)
		case TrapPUTS:
			address := c.Reg[0]
			//	log.Println(fmt.Sprintf("0x%04x: PUTS Address: 0x%04x %c", c.PC, address, c.ReadMemory(address)))

			//	for i, r := range slab {
			//		if r != unicode.ReplacementChar {
			//			// Print the position of the rune's start byte in the string
			//			// and the value of the rune as a Unicode code point.
			//			fmt.Printf("\t%2d: %c\n", i, rune(r))
			//		} else {
			//			// A byte that doesn't form a part of a rune was detected within the string.
			//			// Print the position of the byte in the string
			//			// and the integer value of the byte in hexadecimal.
			//			fmt.Printf("\t%2d: %02X    (not a rune)\n", i, slab[i])
			//		}
			//	}

			//for _, num := range slab {
			//	fmt.Println("sum:", rune(num))
			//}

			var chr uint16
			var i uint16 = 1
			//var byteArr []byte
			for ok := true; ok; ok = (chr != 0x0) {
				chr = c.Memory[address+i] & 0xFFFF
				//chr += c.Memory[address+i+1]

				//log.Println("char is: %d", chr)
				//log.Println(fmt.Sprintf("Address: 0x%04x", address+i))

				//fmt.Printf("\t%2d: %c\n", i, rune(r))
				//chr = slab[i]
				//  c.ReadMemory(address + i)
				//	chr = c.Memory[address+i]
				//	fmt.Println("loop called chr: %d", chr)
				//fmt.Printf("XXX:%c", rune(chr))
				//fmt.Printf("\t%2d: %c\n", i, rune(r))
				//printBytes(rune(chr))
				//fmt.Printf("\t%2d: %U %c\n", i, chr, rune(chr))
				fmt.Printf("%c", rune(chr))

				//fmt.Printf("Byte %3d: %3d %c\n", i, chr, chr)
				//fmt.Printf("\t%2d: %U %c\n", i, chr, chr)

				//fmt.Sprintln("%c %U", chr, chr)

				i++
			}

			fmt.Printf("\n")
			//fmt.Println("Terminating VM")
			//os.Exit(1)
			//for c > 0 {
			//	chr := rune(c)
			//	fmt.Printf("%c", chr)
			//}
			//for c {
		//fmt.Printf("%s", c)
		//	c++
		//}
		case TrapHALT:
			log.Println("HALT")
			os.Exit(1)
		default:
			log.Fatalf("Trap code not implemented: 0x%04X", instr)
		}
	case OpRES:
	case OpRTI:
	default:
		log.Fatalf("Bad Op Code received: %d", op)
	}

	// increment the program counter
	//log.Println(fmt.Sprintf("Setting PC to 0x%04x", pc))
	c.PC = pc
	return
}

func printBytes(s string) {
	fmt.Println("printBytes:")
	sbytes := []byte(s)
	for i, b := range sbytes {
		// Print the position of the byte in the string
		// and the integer value of the byte in hexadecimal.
		fmt.Printf("\t%2d: %2X\n", i, b)
	}
}

func (c *CPU) SetCC(data uint16) {
	c.CondRegister.N = isNegative(data)
	c.CondRegister.Z = isZero(data)
	c.CondRegister.P = isPositive(data)
}

func isPositive(data uint16) bool {
	return int16(data) > 0
}

func isZero(data uint16) bool {
	return data == 0
}

func isNegative(data uint16) bool {
	return int16(data) < 0
}

func extract1C(inst uint16, hi, lo int) uint16 {
	//fmt.Printf("Inst %04x %d %d ", inst, hi, lo)
	if hi >= 16 || hi < 0 || lo >= 16 || lo < 0 {
		fmt.Println("Argument out of bounds")
	}

	//Build mask
	mask := uint16(0)
	for i := 0; i <= hi-lo; i++ {
		mask = mask << 1
		mask |= 0x0001
	}
	for i := 0; i < lo; i++ {
		mask = mask << 1
	}
	//fmt.Printf("Mask %04x ", mask)

	//Apply mask
	field := inst & mask

	//Shift field down
	field = field >> uint(lo)

	//fmt.Printf("Field %04x\n", field)
	return field
}

func extract2C(inst uint16, hi, lo int) uint16 {
	field := extract1C(inst, hi, lo)

	//fmt.Printf("Field %016b ", field)
	if extract1C(field, hi, hi) == 1 {
		//Build sign extension

		mask := uint16(0)
		for i := 0; i <= 15-hi; i++ {
			mask = mask << 1
			mask |= 0x0001
		}
		mask = mask << uint(hi)
		field = inst | mask

	}
	//fmt.Printf("Field %016b\n", field)

	return field
}
