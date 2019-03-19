package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
)

func main() {
	// we need a parallel OS thread to avoid audio stuttering
	//runtime.GOMAXPROCS(2)

	// we need to keep OpenGL calls on a single thread
	//runtime.LockOSThread()

	log.Printf("Starting LC3-VM")

	// load the ROM file
	path := "rom/2048.obj"
	log.Printf("Loading Program: %s", path)

	// read the rom file into a buffer
	mem, err := RetrieveROM(path)
	if err != nil {
		//return nil, err
		panic(err)
	}

	//err = termbox.Init()
	//if err != nil {
	//	panic(err)
	//}
	//defer termbox.Close()

	//eventQueue := make(chan termbox.Event)
	//go func() {
	//	for {
	//		eventQueue <- termbox.PollEvent()
	//	}
	//}()

	// init the CPU
	fmt.Println("Boot VM")
	cpu := NewCPU()
	cpu.Memory = mem
	cpu.Reset()

	cpu.Run()
	fmt.Println("Exiting")
}

func RetrieveROM(filename string) ([65536]uint16, error) {
	m := [65536]uint16{}

	file, err := os.Open(filename)

	if err != nil {
		return m, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return m, statsErr
	}

	// Read origin
	// The first 16 bits of the program file specify the address in memory where the
	// program should start. This address is called the origin.
	var origin uint16

	headerBytes := make([]byte, 2)
	_, err = file.Read(headerBytes)
	if err != nil {
		return m, err
	}

	headerBuffer := bytes.NewBuffer(headerBytes)
	// LC-3 programs are big-endian, but most of the modern computers we use are little endian
	err = binary.Read(headerBuffer, binary.BigEndian, &origin)
	if err != nil {
		return m, err
	}

	log.Printf("Origin memory location: 0x%04X", origin)
	var size int64 = stats.Size()
	byteArr := make([]byte, size)

	log.Printf("Creating memory buffer: %d bytes", size)

	_, err = file.Read(byteArr)
	if err != nil {
		return m, err
	}

	buffer := bytes.NewBuffer(byteArr)

	//bufr := bufio.NewReader(file)
	//_, err = bufr.Read(byteArr)

	//reader := bytes.NewReader(byteArr)

	for i := origin; i < math.MaxUint16; i++ {
		var val uint16
		binary.Read(buffer, binary.BigEndian, &val)
		//fmt.Println("in loop", val, "writing to:", i)
		m[i] = val
	}

	return m, err
}

func readStdin(out chan string, in chan bool) {
	//no buffering
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	//no visible output
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 1)
	for {
		select {
		case <-in:
			return
		default:
			os.Stdin.Read(b)
			fmt.Printf(">>> %v: ", b)
			out <- string(b)
		}
	}
}
