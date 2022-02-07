package main

import (
	"fmt"
	"hzip/compression"
	"hzip/input"
	"hzip/output"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: hzip <output_file> <input_file>")
		fmt.Println("[FATAL] Invalid arguments")
		os.Exit(1)
	}

	compressor := compression.CreateCompressor()

	compressor.SetOutput(output.FileOutput{
		Filename: output.GetOutputFilename(os.Args[1]),
	})
	fmt.Println("[INFO] Collecting input files")
	for _, input_filename := range os.Args[2:] {
		objs, err := input.ExpandInput(input_filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("[FATAL] Input collection failed")
			os.Exit(1)
		}
		for _, input_obj := range objs {
			compressor.AddInput(input_obj)
		}
	}

	fmt.Println("[INFO] Compressing")
	err := compressor.Process()
	if err != nil {
		fmt.Println(err)
		fmt.Println("[FATAL] Compression failed")
		os.Exit(1)
	}

	fmt.Println("[INFO] Dumping archive")
	err = compressor.Dump()
	if err != nil {
		fmt.Println(err)
		fmt.Println("[FATAL] Dump failed")
		os.Exit(1)
	}
}
