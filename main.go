package main

import (
	"flag"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
	"os"
	"transcoder/amrwb"
)

func readWavPCM16(path string) ([]int16, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := wav.NewDecoder(f)
	buf := make([]int16, 0)
	for {
		sample, err := r.FullPCMBuffer()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		buf = append(buf, int16(sample.Data[0]))
	}
	return buf, nil
}

func writeWavPCM16(path string, samples []int16) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := wav.NewEncoder(f, 16000, 16, 1, 1)
	data := &audio.IntBuffer{
		Data:           make([]int, len(samples)),
		Format:         &audio.Format{SampleRate: 16000, NumChannels: 1},
		SourceBitDepth: 16,
	}
	for i, s := range samples {
		data.Data[i] = int(s)
	}
	return w.Write(data)
}

func main() {
	inWav := flag.String("in", "input.wav", "Path to input wav file")
	outWav := flag.String("out", "output.wav", "Path to output wav file")
	mode := flag.Int("mode", 7, "AMR-WB mode (0-8)")
	flag.Parse()

	pcm, err := readWavPCM16(*inWav)
	if err != nil {
		fmt.Println("Read WAV error:", err)
		return
	}

	codec := amrwb.NewAMRWBCodec(int16(*mode), false)
	defer codec.Close()

	var reconstructed []int16
	for i := 0; i+320 <= len(pcm); i += 320 {
		frame := pcm[i : i+320]
		bitstream, err := codec.Encode(frame)
		if err != nil {
			fmt.Println("Encode error:", err)
			return
		}
		signal, err := codec.Decode(bitstream)
		if err != nil {
			fmt.Println("Decode error:", err)
			return
		}
		reconstructed = append(reconstructed, signal...)
	}

	err = writeWavPCM16(*outWav, reconstructed)
	if err != nil {
		fmt.Println("Write WAV error:", err)
		return
	}

	fmt.Println("AMR-WB encode/decode complete.")
}
