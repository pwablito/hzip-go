package main

import (
	"fmt"
	"hzip/src/compression"
	"hzip/src/input"
	"hzip/src/output"
	"os"
)

func main() {
	// Check if compressing or decompressing
	if len(os.Args) < 2 {
		fmt.Println("[FATAL] Must supply a command")
		os.Exit(1)
	}
	if os.Args[1] == "c" || os.Args[1] == "compress" {
		if len(os.Args) < 4 {
			fmt.Println("[FATAL] Arguments to compress missing")
			os.Exit(1)
		}
		outputFilename := os.Args[2]
		inputs := os.Args[3:]

		compressor := compression.CreateCompressor()

		compressor.SetOutput(&output.FileOutput{
			Filename: output.GetOutputFilename(outputFilename),
			Mode:     0666,
		})
		fmt.Println("[INFO] Collecting input files")
		for _, inputFilename := range inputs {
			objs, err := input.ExpandInput(inputFilename)
			// TODO make sure all inputs are in a subdirectory of the current directory
			// TODO if there is only one file, it can be anywhere and will expand to just the object (no dirs)
			if err != nil {
				fmt.Println(err)
				fmt.Println("[FATAL] Input collection failed")
				os.Exit(1)
			}
			for _, inputObj := range objs {
				compressor.AddInput(inputObj)
			}
		}
		// TODO Remove duplicate inputs

		fmt.Println("[INFO] Compressing")
		err := compressor.GenerateScheme()
		if err != nil {
			fmt.Println(err)
			fmt.Println("[FATAL] Compression scheme generation failed")
			os.Exit(1)
		}

		fmt.Println("[INFO] Compressing to archive")
		err = compressor.CompressToOutput()
		if err != nil {
			fmt.Println(err)
			fmt.Println("[FATAL] Dump failed")
			os.Exit(1)
		}
	} else if os.Args[1] == "d" || os.Args[1] == "decompress" {
		if len(os.Args) < 3 {
			fmt.Println("[FATAL] Must supply an archive as an argument")
			os.Exit(1)
		}
		inputFilename := os.Args[2]
		decompressor := compression.CreateDecompressor(inputFilename)
		err := decompressor.ReadMeta()
		if err != nil {
			fmt.Println(err)
			fmt.Println("[FATAL] Failed to read metadata from archive")
			os.Exit(1)
		}
		err = decompressor.Decompress()
		if err != nil {
			fmt.Println(err)
			fmt.Println("[FATAL] Failed to decompress")
			os.Exit(1)
		}
	} else {
		fmt.Println("[FATAL] Invalid command")
		os.Exit(1)
	}
}
