package main

import (
	"fmt"
	"testing"
)

func TestCPUAddInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x3000] = 0x1261 // ADD R1, R1
	m[0x3001] = 0x14A1 // ADD R2, R2, #1

	cpu := NewCPU()
	cpu.Memory = m
	cpu.Reset()

	dr := 1
	cpu.Reg[1] = 9
	cpu.Step()
	if cpu.Reg[dr] != 10 {
		t.Errorf("c.Reg[1] %v expected %v", cpu.Reg[dr], 10)
	}

	dr2 := 2
	cpu.Step()
	if cpu.Reg[dr2] != 1 {
		t.Errorf("c.Reg[2] %v expected %v", cpu.Reg[dr2], 1)
	}
}

func TestCPUBrzInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x0454] = 0x6200 // LDR R1, R0, #0
	m[0x0455] = 0x0405 // BRz x045B
	m[0x0456] = 0xA409 // LDI R2, x0460
	m[0x0457] = 0x07FE // BRzp x0456
	m[0x0458] = 0xB208 // STI R1, x0461
	m[0x0459] = 0x1021 // ADD R0, R0, #1
	m[0x045A] = 0x0FF9 // BRnzp x0454

	cpu := initCPU(m)
	//dr := 1
	cpu.Reg[0] = 0x3080
	cpu.Reg[5] = 0x3017
	cpu.Reg[6] = 0x4000
	cpu.Reg[7] = 0x3004
	cpu.PC = 0x0454

	cpu.Step()
	dumpCPUState(t, cpu)
	t.FailNow()

	// execute the 7 instructions above
	for i := 0; i < 7; i++ {
		fmt.Println("step")
		cpu.Step()
	}

	dumpCPUState(t, cpu)
	t.FailNow()

	// We should of jumped from 0x045A back to 0x0454
	if cpu.PC != 0x0454 {
		dumpCPUState(t, cpu)
		t.Errorf("c.PC 0x%04x expected 0x%04x", cpu.PC, 0x0454)
	}
}

func TestCPUJmpInstr(t *testing.T) {
	//m := [65536]uint16{}
	//m[0x3000] = 0x5260 // AND R1, R1, #0

	//cpu := initCPU(m)
	//dr := 1
	//cpu.Reg[1] = 2
	//cpu.Step()
	//if cpu.Reg[dr] != 12825 {
	//	t.Errorf("c.Reg[1] %v expected %v", cpu.Reg[dr], 12825)
	//}
}

func TestCPULeaInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x31E7] = 0xE031 // LEA R0, x3219
	m[0x3219] = 0x0000

	cpu := NewCPU()
	cpu.Memory = m
	cpu.Reset()
	cpu.PC = 0x31E7

	dr := 0
	cpu.Step()
	if cpu.Reg[dr] != 12825 {
		t.Errorf("c.Reg[0] %v expected %v", cpu.Reg[dr], 12825)
	}
}

func TestCPUAndInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x3000] = 0x5260 // AND R1, R1, #0

	cpu := initCPU(m)
	dr := 1
	cpu.Reg[1] = 2
	cpu.Step()
	if cpu.Reg[dr] != 12825 {
		t.Errorf("c.Reg[1] %v expected %v", cpu.Reg[dr], 12825)
	}
}

func TestCPUJsrInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x32C6] = 0x486D // JSR 0x3334
	m[0x3001] = 0x1020 // ADD R0, R0, #0
	m[0x3334] = 0x1020 // ADD R0, R0, #0

	cpu := initCPU(m)
	cpu.PC = 0x32C6
	dumpCPUState(t, cpu)
	cpu.Step()

	if cpu.PC != 0x3334 {
		dumpCPUState(t, cpu)
		t.Errorf("c.PC 0x%04x expected 0x%04x", cpu.PC, 0x3334)
	}

	// TODO - test RET/JSRR
}

func TestCPULdiInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x0456] = 0xA409 // LDI R2, x0460
	m[0xFE02] = 0x8000 // RTI

	cpu := initCPU(m)
	dr := 2
	cpu.Reg[2] = 0x8000
	cpu.PC = 0x0456
	cpu.Step()
	if cpu.CondRegister.N != true {
		t.Error("c.CondRegister.N should be true")
	}
	if cpu.Reg[dr] != 0x8000 {
		t.Errorf("c.Reg[2] 0x%04x expected 0x%04x", cpu.Reg[dr], 0x8000)
	}
}

func TestCPULdrInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x0454] = 0x6200 // LDR R1, R0, #0

	cpu := initCPU(m)
	dr := 1
	cpu.Reg[0] = 0x3080
	cpu.PC = 0x0454
	cpu.Step()
	if cpu.CondRegister.P != true {
		t.Error("c.CondRegister.P should be true")
	}
	if cpu.Reg[dr] != 0x0043 {
		t.Errorf("c.Reg[1] 0x%04x expected 0x%04x", cpu.Reg[dr], 0x0043)
	}
}

// TODO - finish
func TestCPULdInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x3000] = 0x1261 // ADD R1, R1
}

func initCPU(m [65536]uint16) *CPU {
	cpu := NewCPU()
	cpu.Memory = m
	cpu.Reset()
	return cpu
}

// This method dumps the current state of the CPU in the event of a test failure.
func dumpCPUState(t *testing.T, c *CPU) {
	t.Logf("======== CPU STATE ===========")
	t.Logf("Register 0: 0x%04X", c.Reg[0])
	t.Logf("Register 1: 0x%04X", c.Reg[1])
	t.Logf("Register 2: 0x%04X", c.Reg[2])
	t.Logf("Register 3: 0x%04X", c.Reg[3])
	t.Logf("Register 4: 0x%04X", c.Reg[4])
	t.Logf("Register 5: 0x%04X", c.Reg[5])
	t.Logf("Register 6: 0x%04X", c.Reg[6])
	t.Logf("Register 7: 0x%04X", c.Reg[7])
	t.Logf("N: %v", c.CondRegister.N)
	t.Logf("Z: %v", c.CondRegister.Z)
	t.Logf("P: %v", c.CondRegister.P)
	t.Logf("PC: 0x%04X", c.PC)
	//t.Logf("Inst: 0x%04X Op: %d", instr, op)
}
