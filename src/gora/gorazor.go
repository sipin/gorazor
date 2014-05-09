package main

// considering import fsnotify

import (
	"fmt"
	"os"
	"flag"
	_ "gorazor"
	"errors"
	"path/filepath"
)

const (
	go_extension = ".go"
	gz_extension = ".gohtml"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: in-dir out-dir\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func exits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil}
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func ProcessFile(infile string, indir string, outdir string) (err error) {
	fmt.Printf("processing: %s %s %s\n", infile, indir, outdir)
	return nil
}

func ProcessFolder(indir string, outdir string) (err error) {
	if ok, err := exits(indir); ok == false {
		return errors.New("input directory is not exsits")
	} else {
		if err != nil { return err}
	}
	ok, err := exits(outdir)
	if ok {  // remove original directory
		os.RemoveAll(outdir)
	}
	os.MkdirAll(outdir, 0777)
	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ProcessFile(path, indir, outdir)
		}
		return nil
	}
	err = filepath.Walk(indir, visit)
	return nil
}


func main() {
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()

	in, out := args[0], args[1]
	err := ProcessFolder(in, out)
	if err != nil {
		fmt.Println(err)
	}
}
