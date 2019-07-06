hello:
	./gorazor ./docs/hello/tpl ./docs/hello/src/tpl

examples:
	./gorazor ./examples/tpl ./examples/gen

bench:
	./gorazor ./tests/tpl ./tests/tpl