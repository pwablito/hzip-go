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

	output_obj := output.FileOutput{
		Filename: os.Args[1],
	}

	inputs := make([]input.Input, 0)

	for _, input_filename := range os.Args[2:] {
		inputs = append(inputs, input.FileInput{
			Filename: input_filename,
		})
	}

	compressor := compression.Compressor{
		Output: output_obj,
		Inputs: inputs,
	}

	compressor.Compress()
	compressor.Dump()
}
