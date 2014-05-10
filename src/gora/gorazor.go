package main

// considering import fsnotify

import (
	"fmt"
	"os"
	"flag"
	"gorazor"
)


func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: Specify template file or directory\n")
	fmt.Fprintf(os.Stderr, "      -f=\"\"      : Template file path\n")
        fmt.Fprintf(os.Stderr, "      -o=\"\"      : Output file path\n")
        fmt.Fprintf(os.Stderr, "      -indir=\"\"  : Template directory path\n")
        fmt.Fprintf(os.Stderr, "      -outdir=\"\" : Output directory path\n")
	os.Exit(0)
}

func main() {

	var InDir, OutDir, InFile, OutFile string
	flag.StringVar(&InDir, "indir", "", "Template directory path")
	flag.StringVar(&OutDir, "outdir", "", "Output directory path")
	flag.StringVar(&InFile, "f", "", "Template file path")
	flag.StringVar(&OutFile, "o", "", "Output file path")
	flag.Usage = Usage

	flag.Parse()

	if InDir != "" && OutDir != "" {
		err := gorazor.GenFolder(InDir, OutDir)
		if err != nil {
			fmt.Println(err)
		}
	} else if InFile != "" && OutFile != ""  {
		fmt.Printf("processing: %s %s\n", InFile, OutFile)
		gorazor.GenFile(InFile, OutFile)
	} else {
		flag.Usage()
	}
}
