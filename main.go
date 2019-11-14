package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sipin/gorazor/pkg/razorcore"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gorazor [-debug] <input dir or file> <output dir or file>\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	isDebug := flag.Bool("debug", false, "use debug mode")
	noLine := flag.Bool("noline", false, "skip line number hint in generated code")
	version := flag.Bool("version", false, "show gorazor version info")
	quick := flag.Bool("q", false, "enable quick mode; skip template render optimzation")
	namespacePrefix := flag.String("prefix", "", "tpl namespace prefix")
	nameNotChange := flag.Bool("nameNotChange", false, "do not change name of the template")

	flag.Parse()

	if *version {
		fmt.Println("gorazor version: " + razorcore.VERSION)
		os.Exit(0)
	}

	options := razorcore.Option{}

	options.IsDebug = *isDebug
	options.NameNotChange = *nameNotChange
	options.NoLineNumber = *noLine

	if len(flag.Args()) != 2 {
		flag.Usage()
	}

	input, output := flag.Arg(0), flag.Arg(1)
	stat, err := os.Stat(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	razorcore.TemplateNamespacePrefix = *namespacePrefix
	razorcore.QuickMode = *quick

	if stat.IsDir() {
		fmt.Printf("gorazor processing dir: %s -> %s\n", input, output)
		err := razorcore.GenFolder(input, output, options)
		if err != nil {
			fmt.Println(err)
		}
	} else if stat.Mode().IsRegular() {
		fmt.Printf("gorazor processing file: %s -> %s\n", input, output)
		razorcore.GenFile(input, output, options)
	}
}
