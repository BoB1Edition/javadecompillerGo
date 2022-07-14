package main

import (
	"flag"
	"fmt"
	"localhost/javadecompiler/decompiler"
	"log"
	"os"
	"path"
	"strings"
)

var (
	output    *string
	outputDir *string
	fName     string
)

func init() {
	if len(os.Args) > 0 {
		fName = os.Args[len(os.Args)-1]
		output = new(string)
		outputDir = new(string)
		defOutput := fName[0:len(fName)-len(path.Ext(fName))] + ".java"
		flag.StringVar(output, "output", defOutput, "output file")
		flag.StringVar(outputDir, "outputdir", "./", "path to outputdir, only war or jar file")
		flag.Parse()
	} else {
		fmt.Printf("please use help\n\tRequired argument not specified FILENAME")
	}
}

func main() {
	if fName == "" {
		fmt.Print("input file not specified")
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	switch strings.ToLower(path.Ext(fName))[1:] {
	case "war":
		fmt.Print("TODO: war file")
	case "jar":
		fmt.Print("TODO: jar file")
	case "class":
		d := decompiler.New(fName)
		if err := d.ParseFile(); err != nil {
			log.Panic(err)
		}
		d.WriteFile(*output)
	}
}
