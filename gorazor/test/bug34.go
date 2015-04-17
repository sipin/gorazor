package cases

import (
	"bytes"
)

func Bug34() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("value=\\\"<?= h(aabasdf\\Admin\\Document::$asdf) ?>\\\"/>\\n")

	return _buffer.String()
}
