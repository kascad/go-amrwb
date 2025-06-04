// Package amrwb implements a partial translation of AMR-WB decoder logic
package decoder

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	LFrame16k   = 320
	NbSerialMax = 61
	AmrWbMagic  = "#!AMR-WB\n"
)

type DecoderState struct {
	// Placeholder for actual decoder internal state
}

func D_IF_Init() *DecoderState {
	// Initialize decoder state (placeholder)
	return &DecoderState{}
}

func D_IF_Exit(st *DecoderState) {
	// Cleanup decoder state (placeholder)
}

func D_IF_Decode(st *DecoderState, serial []byte, synth []int16) int {
	// Placeholder decode logic, returns dummy mode
	for i := range synth {
		synth[i] = 0
	}
	return 0 // dummy mode
}

func RunDecoder(inputPath, outputPath string) error {
	fSerial, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer fSerial.Close()

	fSynth, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer fSynth.Close()

	fmt.Println("Input bitstream file:", inputPath)
	fmt.Println("Synthesis speech file:", outputPath)

	st := D_IF_Init()
	defer D_IF_Exit(st)

	// Check and skip AMR-WB magic number if present
	magic := make([]byte, len(AmrWbMagic))
	if _, err := io.ReadFull(fSerial, magic); err != nil {
		return fmt.Errorf("error reading magic number: %w", err)
	}
	if string(magic) != AmrWbMagic {
		return fmt.Errorf("invalid magic number: %q", string(magic))
	}

	serial := make([]byte, NbSerialMax)
	synth := make([]int16, LFrame16k)
	frame := 0

	for {
		// Read 1st byte (frame header)
		n, err := fSerial.Read(serial[:1])
		if err == io.EOF {
			break
		}
		if err != nil || n != 1 {
			return fmt.Errorf("error reading frame header: %w", err)
		}

		// Read the rest of the frame (assuming max size)
		mode := serial[0] >> 3 & 0x0F // dummy mode extraction
		frameSize := int(blockSize[mode])

		if _, err := io.ReadFull(fSerial, serial[1:frameSize]); err != nil {
			return fmt.Errorf("error reading frame body: %w", err)
		}

		D_IF_Decode(st, serial[:frameSize], synth)

		// Write synthesized frame to output
		if err := binary.Write(fSynth, binary.LittleEndian, synth); err != nil {
			return fmt.Errorf("error writing synthesized audio: %w", err)
		}

		frame++
	}

	fmt.Printf("\n%d frames processed.\n", frame)
	return nil
}

// blockSize corresponds to frame sizes per mode (dummy placeholder)
var blockSize = [...]uint8{
	17, 23, 32, 36, 40, 46, 50, 58, 60, 5, 0, 0, 0, 0, 0, 0,
}
