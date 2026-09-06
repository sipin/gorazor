package razorcore

//------------------------------ API ------------------------------

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	goExtension = ".go"
	gzExtension = ".gohtml"
)

// GenFile generate from input to output file,
// gofmt will trigger an error if it fails.
func GenFile(ctx context.Context, input string, output string, options Option) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	outdir := filepath.Dir(output)
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}
	if options.LayoutCache == nil {
		options.LayoutCache = NewLayoutCache()
	}
	return generate(ctx, input, output, options)
}

// GenFolder generate from directory to directory, Find all the files with extension
// of .gohtml and generate it into target dir.
func GenFolder(ctx context.Context, indir string, outdir string, options Option) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !exists(indir) {
		return errors.New("input directory does not exsit")
	}

	if options.LayoutCache == nil {
		options.LayoutCache = NewLayoutCache()
	}

	//Make it
	if !exists(outdir) {
		os.MkdirAll(outdir, 0775)
	}

	incdirAbs, _ := filepath.Abs(indir)
	outdirAbs, _ := filepath.Abs(outdir)

	paths := []string{}

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//Just do file with exstension .gohtml
			if !strings.HasSuffix(path, gzExtension) {
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

	var errMutex sync.Mutex
	var firstErr error

	fun := func(path string, res chan<- string) {
		defer func() {
			if r := recover(); r != nil {
				errMutex.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("error processing %s: %v", path, r)
				}
				errMutex.Unlock()
				res <- fmt.Sprintf("error: %s: %v", path, r)
			}
		}()

		if err := ctx.Err(); err != nil {
			errMutex.Lock()
			if firstErr == nil {
				firstErr = err
			}
			errMutex.Unlock()
			res <- fmt.Sprintf("cancelled: %s", path)
			return
		}

		rel, relErr := filepath.Rel(indir, path)
		var output string
		if relErr == nil {
			output = filepath.Join(outdir, rel)
		} else {
			input, _ := filepath.Abs(path)
			output = strings.Replace(input, incdirAbs, outdirAbs, 1)
		}
		output = strings.ReplaceAll(output, gzExtension, goExtension)
		genErr := GenFile(ctx, path, output, options)
		if genErr != nil {
			errMutex.Lock()
			if firstErr == nil {
				firstErr = genErr
			}
			errMutex.Unlock()
		}
		res <- fmt.Sprintf("%s -> %s", path, output)
	}

	err = filepath.Walk(indir, visit)
	if err != nil {
		return err
	}
	result := make(chan string, len(paths))
	jobs := make(chan string, len(paths))
	for _, p := range paths {
		jobs <- p
	}
	close(jobs)

	numWorkers := runtime.NumCPU()
	if numWorkers > len(paths) {
		numWorkers = len(paths)
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				fun(p, result)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for res := range result {
		fmt.Println(res)
	}

	return firstErr
}
