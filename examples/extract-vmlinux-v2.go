package main

/*
	Original source: https://github.com/Caesurus/extract-vmlinux-v2
*/

import (
	"fmt"
	vmlinux "github.com/soxfmr/extract-vmlinux-v2"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser(appName, "A more robust vmlinux extractor")
	var inputFile = parser.String("f", "file", &argparse.Options{Help: "Input file to be processed"})
	var outputFile = parser.String("f", "output", &argparse.Options{Help: "Output file to be stored the vmlinux"})
	var printVersion = parser.Flag("V", "version", &argparse.Options{Help: "Print version info"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *printVersion {
		fmt.Printf("%s: %d.%d.%d\n", appName, versionMajor, versionMinor, versionPatch)
		os.Exit(0)
	}

	if len(*inputFile) > 0 {
		kernelFile, err := os.Open(*inputFile)
		if err != nil {
			log.Fatalf("couldn't open the input file: %s", err)
		}

		var extractedFile *os.File
		if len(*outputFile) > 0 {
			extractedFile, err = os.Open(*outputFile)
		} else {
			extractedFile, err = os.CreateTemp("", "vmlinux")
		}
		if err != nil {
			log.Fatalf("couldn't open the output file: %s", err)
		}

		if err := vmlinux.ExtractTo(kernelFile, extractedFile); err != nil {
			log.Fatalf("couldn't extract the kernel image: %s", err)
		}
	}

}
