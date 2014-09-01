package gorazor

//------------------------------ API ------------------------------

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	go_extension = ".go"
	gz_extension = ".gohtml"
)

// Generate from input to output file,
// gofmt will trigger an error if it fails.
func GenFile(input string, output string, options Option) error {
	outdir := filepath.Dir(output)
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}
	return generate(input, output, options)
}

// Generate from directory to directory, Find all the files with extension
// of .gohtml and generate it into target dir.
func GenFolder(indir string, outdir string, options Option) (err error) {
	if !exists(indir) {
		return errors.New("Input directory does not exsits")
	} else {
		if err != nil {
			return err
		}
	}
	//Make it
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}

	incdir_abs, _ := filepath.Abs(indir)
	outdir_abs, _ := filepath.Abs(outdir)

	paths := []string{}

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//Just do file with exstension .gohtml
			if !strings.HasSuffix(path, gz_extension) {
				return nil
			}
			filename := filepath.Base(path)
			if strings.HasPrefix(filename, ".#") {
				return nil
			}
			paths = append(paths, path)
		}
		return nil
	}

	fun := func(path string, res chan<- string) {
		//adjust with the abs path, so that we keep the same directory hierarchy
		input, _ := filepath.Abs(path)
		output := strings.Replace(input, incdir_abs, outdir_abs, 1)
		output = strings.Replace(output, gz_extension, go_extension, -1)
		err := GenFile(path, output, options)
		if err != nil {
			res <- fmt.Sprintf("%s -> %s", path, output)
			os.Exit(2)
		}
		res <- fmt.Sprintf("%s -> %s", path, output)
	}

	err = filepath.Walk(indir, visit)
	runtime.GOMAXPROCS(runtime.NumCPU())
	result := make(chan string, len(paths))

	for w := 0; w < len(paths); w++ {
		go fun(paths[w], result)
	}
	for i := 0; i < len(paths); i++ {
		<-result
	}

	if options["Watch"] != nil {
		watchDir(incdir_abs, outdir_abs, options)
	}
	return
}
