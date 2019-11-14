# Backlog

- [ ] html compact mode option
- [ ] Return error during rendering?
- [ ] Better error msg during codegen
- [ ] Performance Optimize
  - [ ] Auto convert helper func to writer
  - [ ] Unsafe write?
- [ ] VS Code plugin
- [ ] Quick Start guide
- [ ] Dynamic compile?
- [ ] Support webassembly
- [ ] Performance Optimize
  - [ ] Auto convert helper func to writer
  - [ ] Unsafe write?

# v1.2.2
- [X] Line number support
  - [X] option
- [X] Refactor
  - [X] cmd options into struct
  - [X] refactor lexer, keep line number

# v1.2.1
- [X] Add version
  - [X] codegen header
  - [X] version output
- [X] Add namespace prefix fix support
- [X] Define new layout syntax
  - [X] Update readme
  - [X] Decide layout import path discovery
  - [X] Implement `layout := XXX`
  - [X] Implement `islayout := true`
- [X] Improve test
  - [X] Add bench test with quicktemplate
  - [X] More test case files
  - [X] Make code coverage 90%+
  - [X] Test cases for QuickMode
- [X] Proper tagging for version
- [X] Performance Optimize
  - [X] Setup benchmark
  - [X] zero alloc
- [X] Refactor
  - [X] Put utils in independent namespace
