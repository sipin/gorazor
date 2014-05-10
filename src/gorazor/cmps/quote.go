package test

import (
"bytes"
)

func Quote() (string) {
var _buffer bytes.Buffer
_buffer.WriteString("<html>'text\'</html>\n")
return _buffer.String()
}
