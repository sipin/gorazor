package main

// considering import fsnotify

import (
	"flag"
	"fmt"
	"os"

	"github.com/sipin/gorazor/gorazor"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: gorazor <input dir or file> <output dir or file> [-debug]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	options := gorazor.Option{}

	if flag.NArg() == 3 && flag.Arg(2) == "-debug" {
		options["Debug"] = true
	}
	if flag.NArg() >= 2 {
		arg1, arg2 := flag.Arg(0), flag.Arg(1)
		stat, err := os.Stat(arg1)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if stat.IsDir() {
			fmt.Printf("Processing dir: %s %s\n", arg1, arg2)
			err := gorazor.GenFolder(arg1, arg2, options)
			if err != nil {
				fmt.Println(err)
			}
		} else if stat.Mode().IsRegular() {
			fmt.Printf("Processing file: %s %s\n", arg1, arg2)
			gorazor.GenFile(arg1, arg2, options)
		} else {
			flag.Usage()
		}
	}
}
