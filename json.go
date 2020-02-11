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

func newJSON() *json {
	js := jsonPool.Get().(*json)
	js.buf = append(js.buf, '{')

	return js
}

func (js *json) put() {
	js.buf = js.buf[:0]
	jsonPool.Put(js)
}

func (js *json) close() {
	js.buf = append(js.buf, '}')
}

func (js *json) addField(key string, value []byte) {
	if key == "" {
		return
	}

	if len(js.buf) > 3 {
		js.buf = append(js.buf, ", "...)
	}

	js.buf = append(js.buf, '"')
	js.buf = append(js.buf, key...)
	js.buf = append(js.buf, '"')
	js.buf = append(js.buf, ": "...)
	js.buf = append(js.buf, '"')
	js.buf = append(js.buf, value...)
	js.buf = append(js.buf, '"')
}
