package main

import (
	"fmt"
	"hzip/src/compression"
	"hzip/src/input"
	"hzip/src/output"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("hzip", "Compression tool")
	output_filename := parser.String("o", "output-filename", &argparse.Options{Required: true, Help: "Output filename"})
	inputs := parser.StringList("i", "input", &argparse.Options{Required: true, Help: "Inputs"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	compressor := compression.CreateCompressor()

	compressor.SetOutput(&output.FileOutput{
		Filename: output.GetOutputFilename(*output_filename),
		Mode:     0666,
	})
	fmt.Println("[INFO] Collecting input files")
	for _, input_filename := range *inputs {
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
	// TODO Removeo duplicate inputs

	fmt.Println("[INFO] Compressing")
	err = compressor.Process()
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
