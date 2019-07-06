hello:
	./gorazor ./docs/hello/tpl ./docs/hello/src/tpl

examples:
	./gorazor -prefix github.com/sipin/gorazor ./examples/tpl ./examples/tpl

bench:
	./gorazor ./tests/tpl ./tests/tpl