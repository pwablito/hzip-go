package main

import (
	"fmt"
	"hzip/compression"
	"hzip/input"
	"hzip/output"
	"os"
)

func print_usage() {
	fmt.Println("Usage: hzip <output_file> <input_file>")
}

func main() {
	if len(os.Args) < 3 {
		print_usage()
		os.Exit(1)
	}

	compressor := compression.CreateCompressor()

	compressor.SetOutput(output.FileOutput{
		Filename: output.GetOutputFilename(os.Args[1]),
	})
	for _, input_filename := range os.Args[2:] {
		compressor.AddInput(input.FileInput{
			Filename: input_filename,
		})
	}

	err := compressor.Compress()
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] compression failed")
		os.Exit(1)
	}
	err = compressor.Dump()
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Dump failed")
		os.Exit(1)
	}
}
