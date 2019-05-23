package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

func processInput(cpu *CPU) (err error) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if cpu.DebugMode {
				log.Println(fmt.Sprintf("Key pressed: %d", ev.Ch))
			}
			cpu.keyBuffer = append(cpu.keyBuffer, ev.Ch)
			switch {
			case ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD:
				instr := cpu.ReadMemory(cpu.PC)
				op := instr >> 12

				// stop the CPU from executing
				cpu.Stop()

				if cpu.DebugMode {
					log.Println("========= DEBUG OUTPUT ====================")
					log.Println(fmt.Sprintf("R0: 0x%04X", cpu.Reg[0]))
					log.Println(fmt.Sprintf("R1: 0x%04X", cpu.Reg[1]))
					log.Println(fmt.Sprintf("R2: 0x%04X", cpu.Reg[2]))
					log.Println(fmt.Sprintf("R3: 0x%04X", cpu.Reg[3]))
					log.Println(fmt.Sprintf("R4: 0x%04X", cpu.Reg[4]))
					log.Println(fmt.Sprintf("R5: 0x%04X", cpu.Reg[5]))
					log.Println(fmt.Sprintf("R6: 0x%04X", cpu.Reg[6]))
					log.Println(fmt.Sprintf("R7: 0x%04X", cpu.Reg[7]))
					log.Println(fmt.Sprintf("PC: 0x%04X", cpu.PC))
					log.Println(fmt.Sprintf("Inst: 0x%04X Op: %d", instr, op))
				}
				return
			}
		default:
			// do nothing
		}
	}
}
