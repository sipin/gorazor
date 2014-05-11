package main

// considering import fsnotify

import (
	"flag"
	"fmt"
	"gorazor"
	"os"
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

	var indir, outdir, infile, outfile string
	debug := false
	flag.StringVar(&indir, "indir", "", "Template directory path")
	flag.StringVar(&outdir, "outdir", "", "Output directory path")
	flag.StringVar(&infile, "f", "", "Template file path")
	flag.StringVar(&outfile, "o", "", "Output file path")
	flag.BoolVar(&debug, "d", false, "Enable debug mode")
	flag.Usage = Usage

	flag.Parse()

	options := gorazor.Option{}
	if debug {
		options["Debug"] = true
	}

	if indir != "" && outdir != "" {
		err := gorazor.GenFolder(indir, outdir)
		if err != nil {
			fmt.Println(err)
		}
	} else if infile != "" && outfile != "" {
		fmt.Printf("processing: %s %s\n", infile, outfile)
		gorazor.GenFile(infile, outfile, options)
	} else {
		flag.Usage()
	}
}
