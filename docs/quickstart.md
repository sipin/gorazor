# Quick Start Guide

Welcome to `gorazor`! This guide will walk you through setting up a simple project, writing your first template, compiling it to Go code, and rendering it inside a standard Go web server.

---

## 1. Installation

First, install the `gorazor` command-line tool on your machine:

```bash
go install github.com/sipin/gorazor@latest
```

Ensure your Go bin directory (`$GOPATH/bin` or `$HOME/go/bin`) is in your system's `PATH`.

---

## 2. Step 1: Write a Template

Create a new directory named `tpl` for your templates, and create a file inside named `tpl/index.gohtml`:

```html
@{
    // The first block is dedicated to declarations
    var name string
    var items []string
}

<!DOCTYPE html>
<html>
<head>
    <title>Gorazor Quick Start</title>
</head>
<body>
    <h1>Hello, @name!</h1>
    
    <h3>Your Shopping List:</h3>
    @if len(items) == 0 {
        <p>No items on your list.</p>
    } else {
        <ul>
        @for _, item := range items {
            <li>@item</li>
        }
        </ul>
    }
</body>
</html>
```

---

## 3. Step 2: Compile the Template

Run the `gorazor` compiler on your `tpl` folder to generate Go source files. Run the following command:

```bash
gorazor tpl tpl
```

This compiles `tpl/index.gohtml` into a native Go source file `tpl/index.go`. Under the hood, this file exports:
* `tpl.Index(name string, items []string) string`: Renders the template to a string.
* `tpl.RenderIndex(_buffer io.StringWriter, name string, items []string)`: Renders directly into an `io.StringWriter` (such as `*strings.Builder`, `*bufio.Writer`, or a custom buffer) for zero-allocation performance. Note: standard `http.ResponseWriter` does not implement `io.StringWriter`, so in HTTP handlers use `io.WriteString(w, tpl.Index(...))` or wrap `w` with `bufio.NewWriter(w)`.

---

## 4. Step 3: Serve the Template in a Go Web Server

Create a `main.go` file in your root folder:

```go
package main

import (
	"io"
	"log"
	"net/http"

	"your_project_module/tpl" // Replace with your actual project's module path
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := "Developer"
		items := []string{"Go", "Gorazor", "Zero-alloc speed"}

		// Render index template to HTTP response writer
		io.WriteString(w, tpl.Index(name, items))

		// Note on Zero-Allocation:
		// Standard http.ResponseWriter does not implement io.StringWriter.
		// For high-throughput zero-allocation streaming, wrap w with a buffered writer:
		//   bw := bufio.NewWriter(w)
		//   tpl.RenderIndex(bw, name, items)
		//   bw.Flush()
	})

	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
```

Start the web server:
```bash
go run main.go
```
Open `http://localhost:8080` in your browser to see your template rendered instantly.

---

## 5. Advanced Optimization & Features

### A. HTML Compact Mode (Minification)
To minify the compiled HTML output, use the `-compact` flag during code generation:

```bash
gorazor -compact tpl tpl
```

* **What it does**: Collapses excessive spaces, tabs, and newlines in HTML elements, stripping whitespaces between consecutive tags while keeping a single space between words.
* **Format Preservation**: It automatically preserves whitespaces inside `<pre>`, `<textarea>`, `<script>`, and `<style>` blocks, ensuring Javascript logic and CSS formats remain fully intact.

### B. Premium Codegen Diagnostics
`gorazor` comes with detailed syntax error reporting. If you make a mistake in a template, the compiler points out the exact file path, failing line number, and a highlighted visual code snippet:

```text
compilation error in tpl/index.gohtml:14: failed to parse imports block: 14:9: missing import path
Context:
---> 14: import = 123
```

### C. Skipping Line Numbers
To omit auto-generated line comment hints in generated code (e.g. `// Line: 12`), use the `-noline` flag:

```bash
gorazor -noline tpl tpl
```
