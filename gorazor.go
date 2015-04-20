package main

// considering import fsnotify

import (
	"flag"
	"fmt"
	"os"

	"github.com/sipin/gorazor/gorazor"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: gorazor [-debug] [-watch] <input dir or file> <output dir or file>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = Usage
	isDebug := flag.Bool("debug", false, "use debug mode")
	isWatch := flag.Bool("watch", false, "use watch mode")
	nameNotChange := flag.Bool("nameNotChange", false, "do not change name of the template")

	flag.Parse()

	options := gorazor.Option{}

	if *isDebug {
		options["Debug"] = *isDebug
	}
	if *isWatch {
		options["Watch"] = *isWatch
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
		fmt.Printf("Gorazor processing dir: %s -> %s\n", input, output)
		err := gorazor.GenFolder(input, output, options)
		if err != nil {
			fmt.Println(err)
		}
	} else if stat.Mode().IsRegular() {
		fmt.Printf("Gorazor processing file: %s -> %s\n", input, output)
		gorazor.GenFile(input, output, options)
	}
}
