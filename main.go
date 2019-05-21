package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	log.Printf("Starting LC3-VM by Rob Morgan")

	// load the ROM file
	path := "rom/2048.obj"
	log.Printf("Loading Program: %s", path)

	// read the rom file into a buffer
	mem, err := RetrieveROM(path)
	if err != nil {
		panic(err)
	}

	// init the CPU
	fmt.Println("Boot VM")
	termbox.Flush()
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
