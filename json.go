package logg

import (
	"sync"
)

type json struct {
	buf []byte
}

var jsonPool = sync.Pool{
	New: func() interface{} {
		return &json{
			buf: make([]byte, 0, 500),
		}
	},
}

func newJson() *json {
	js := jsonPool.Get().(*json)
	js.buf = js.buf[:0]

	js.buf = append(js.buf, '{')

	return js
}

func (js *json) close() {
	if len(js.buf) == 0 {
		js.buf = append(js.buf, "{}"...)
		return
	}

	if js.buf[len(js.buf)-1] != '"' && len(js.buf) != 1 {
		js.buf = append(js.buf, '"')
	}

	js.buf = append(js.buf, '}')
}

func (js *json) put() {
	const maxSize = 1 << 16 // 64KiB
	if cap(js.buf) > maxSize {
		return
	}

	jsonPool.Put(js)
}

func (js *json) addField(key string, dst []byte) []byte {
	if key == "" {
		return dst
	}

	if len(dst) > 3 {
		dst = append(dst, '"')
		dst = append(dst, ", "...)
	}

	dst = append(dst, '"')
	dst = append(dst, key...)
	dst = append(dst, '"')
	dst = append(dst, ": "...)
	dst = append(dst, '"')

	return dst
}
