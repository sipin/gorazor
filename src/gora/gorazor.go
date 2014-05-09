package main

// considering import fsnotify

import (
	"fmt"
	"os"
	"flag"
	"gorazor"
	"errors"
	"strings"
	"os/exec"
	"path/filepath"
	"io/ioutil"
)

const (
	go_extension = ".go"
	gz_extension = ".gohtml"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: Specify template file or directory\n")
	fmt.Fprintf(os.Stderr, "      -f=\"\"      : Template file path\n")
        fmt.Fprintf(os.Stderr, "      -o=\"\"      : Output file path\n")
        fmt.Fprintf(os.Stderr, "      -indir=\"\"  : Template directory path\n")
        fmt.Fprintf(os.Stderr, "      -outdir=\"\" : Output directory path\n")
	os.Exit(0)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil}
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func ProcessFile(file string, indir string, outdir string) (err error) {
	fmt.Printf("processing: %s %s %s\n", file, indir, outdir)
	Options := map[string]interface{}{}
	fabs, _ := filepath.Abs(file)
	iabs, _ := filepath.Abs(indir)
	oabs, _ := filepath.Abs(outdir)
	abs := strings.Replace(fabs, iabs, oabs, 1)
        out := strings.Replace(abs, gz_extension, go_extension, -1)
	dir := filepath.Dir(abs)

	if ok, _ := exists(dir); !ok {
		os.MkdirAll(dir, 0777)
	}

	res, err := gorazor.Generate(file, Options)
	if err != nil {
		panic(err)
	} else {
		err := ioutil.WriteFile(out, []byte(res), 0777)
		if err != nil { panic(err) }
		cmd := exec.Command("gofmt", "-w", out)
		if err := cmd.Run(); err != nil {
			//panic(err)
		}
	}
	return nil
}

func ProcessFolder(indir string, outdir string) (err error) {
	if ok, err := exists(indir); ok == false {
		return errors.New("Input directory does not exsits")
	} else {
		if err != nil { return err}
	}
	if ok, _ := exists(outdir); ok {
		// remove original directory
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
	Options := map[string]interface{}{}

	var InDir, OutDir, InFile, OutFile string
	flag.StringVar(&InDir, "indir", "", "Template directory path")
	flag.StringVar(&OutDir, "outdir", "", "Output directory path")
	flag.StringVar(&InFile, "f", "", "Template file path")
	flag.StringVar(&OutFile, "o", "", "Output file path")
	flag.Usage = Usage

	flag.Parse()

	if InDir != "" && OutDir != "" {
		err := ProcessFolder(InDir, OutDir)
		if err != nil {
			fmt.Println(err)
		}
	} else if InFile != "" && OutFile != ""  {
		fmt.Printf("processing: %s %s\n", InFile, OutFile)
	}
	fmt.Println(Options["debug"])
}
