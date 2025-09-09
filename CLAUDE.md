# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gorazor is a Go port of the Razor view engine, providing a fast template compilation system that converts `.gohtml` templates into Go code. The project emphasizes extreme performance (~20x faster than html/template) by generating compiled Go code instead of using reflection.

## Common Development Commands

### Building and Installation
- `go build` - Build the gorazor binary
- `go install github.com/sipin/gorazor@latest` - Install from source

### Template Generation
- `gorazor template_folder output_folder` - Process entire directory
- `gorazor template_file output_file` - Process single file
- `gorazor -prefix github.com/sipin/gorazor ./examples/tpl ./examples/tpl` - Generate with namespace prefix

### Testing and Benchmarking
- `go test ./...` - Run all tests
- `go test -bench='Benchmark(Razor|RazorQuick|Quick|HTML)Template' -benchmem github.com/sipin/gorazor/tests` - Run performance benchmarks
- `make examples` - Regenerate example templates
- `make bench` - Regenerate benchmark test templates

## Code Architecture

### Core Components

**`main.go`** - CLI entry point that parses flags and delegates to razorcore package

**`pkg/razorcore/`** - Core template processing engine:
- `api.go` - Public API for file/folder generation (`GenFile`, `GenFolder`)
- `lexer.go` - Tokenizes `.gohtml` template syntax 
- `parser.go` - Builds AST from tokens
- `compiler.go` - Generates Go code from AST
- `optimizer.go` - Optimizes generated code for better performance
- `layout.go` - Handles layout and section template functionality

### Template Processing Pipeline

1. **Lexing**: `.gohtml` files → tokens (AT, BRACE_OPEN, IDENTIFIER, etc.)
2. **Parsing**: Tokens → AST with PROGRAM, BLOCK, MARKUP nodes
3. **Compilation**: AST → Go source code with proper escaping and layout support
4. **Optimization**: Generated code optimization for performance (disabled with `-q` flag)

### Key Conventions

**Template Structure:**
- Template files must use `.gohtml` extension
- Generated Go files use template name with first letter capitalized as function name
- Helper templates must be in `helper/` subdirectory  
- Layout templates must be in `layout/` subdirectory
- First code block `@{}` is for declarations (imports, models, layout)

**Package Organization:**
- Template folder name becomes package name in generated code
- Helper functions are under `helper` namespace
- Layout functions are under `layout` namespace

### Template Features

The engine supports Razor-style syntax:
- Variables: `@variable` (auto-escaped), `@raw(variable)` (unescaped)
- Code blocks: `@{ go code here }`
- Flow control: `@if`, `@for` with standard Go syntax
- Sections: `@section name { content }`
- Layouts: Declared via `layout := layout.Base` in first code block

### Performance Features

- Templates compile to native Go code for maximum speed
- Supports quicktemplate-style ByteBuffer optimization for zero allocations
- Optional optimization phase (controlled by `QuickMode` global variable)
- Concurrent processing of multiple templates in folder mode

### Testing Strategy

- Unit tests in `pkg/razorcore/` cover lexing, parsing, AST generation
- Benchmark tests in `tests/` compare performance against html/template and quicktemplate
- Test data in `pkg/razorcore/testdata/` and examples in `examples/`
- Golden file testing approach for template generation validation