package main

import (
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

// TODO - finish
func TestCPULdInstr(t *testing.T) {
	m := [65536]uint16{}
	m[0x3000] = 0x1261 // ADD R1, R1
}
