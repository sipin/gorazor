package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sipin/gorazor/gorazor"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gorazor [-debug] <input dir or file> <output dir or file>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	isDebug := flag.Bool("debug", false, "use debug mode")
	version := flag.Bool("version", false, "show gorazor version info")
	nameNotChange := flag.Bool("nameNotChange", false, "do not change name of the template")

	flag.Parse()

	if *version {
		fmt.Println("gorazor version: " + gorazor.VERSION)
		os.Exit(0)
	}

	options := gorazor.Option{}

	if *isDebug {
		options["Debug"] = *isDebug
	}
	if *nameNotChange {
		options["NameNotChange"] = *nameNotChange
	}

	if len(flag.Args()) != 2 {
		flag.Usage()
	}

	input, output := flag.Arg(0), flag.Arg(1)
	stat, err := os.Stat(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if stat.IsDir() {
		fmt.Printf("gorazor processing dir: %s -> %s\n", input, output)
		err := gorazor.GenFolder(input, output, options)
		if err != nil {
			fmt.Println(err)
		}
	} else if stat.Mode().IsRegular() {
		fmt.Printf("gorazor processing file: %s -> %s\n", input, output)
		gorazor.GenFile(input, output, options)
	}
}
