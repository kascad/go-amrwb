package encoder

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func EncoderMainLoop(inputPath, outputPath string, mode int16, allowDTX bool) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer outFile.Close()

	state := E_MAIN_init(mode, allowDTX)
	defer func() { E_MAIN_exit(&state) }()

	var signal [320]int16
	var serial [61]int16

	_, _ = outFile.WriteString("#!AMR-WB\n")

	for {
		err = binary.Read(inFile, binary.LittleEndian, &signal)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error reading frame: %w", err)
		}

		E_MAIN_encode(state, signal[:], serial[:])

		for _, val := range serial[:] {
			err = binary.Write(outFile, binary.LittleEndian, val)
			if err != nil {
				return fmt.Errorf("error writing serial: %w", err)
			}
		}
	}
	return nil
}
