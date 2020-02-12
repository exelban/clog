package logg

import (
	"bytes"
	"testing"
)

func TestLogg_getJson(t *testing.T) {
	js := newJson()
	if js == nil {
		t.Error("js object is nil")
		t.FailNow()
	}

	if len(js.buf) == 0 {
		t.Error("js buffer is empty, but must contain a opening bracket '{'")
	}
}

func TestLogg_json_close(t *testing.T) {
	tests := map[string]struct {
		buf      []byte
		expected string
	}{
		"empty": {
			buf:      []byte{},
			expected: `{}`,
		},
		"opened": {
			buf:      []byte{'{'},
			expected: `{}`,
		},
		"some value": {
			buf:      []byte(`"key": "value`),
			expected: `"key": "value"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			js := &json{buf: tc.buf}
			js.close()

			if string(js.buf) != tc.expected {
				t.Errorf("wrong buffer. Expected: %s, received: %s", tc.expected, string(js.buf))
			}
		})
	}
}

func TestLogg_json_put(t *testing.T) {
	js := &json{
		buf: make([]byte, 1<<16+1),
	}

	js.close()
	js.put()

	jsonFromPool := jsonPool.Get().(*json)
	if len(jsonFromPool.buf) >= 1<<16 {
		t.Error("json has a buf size maxSize+1, cannot be saved in sync.Pool")
	}
}

func TestLogg_json_addField(t *testing.T) {
	js := &json{}

	js.buf = js.addField("", []byte{})
	if !bytes.Equal(js.buf, []byte{}) {
		t.Error("field with empty key must be missed")
	}

	testString := `"1": "1`
	js.buf = js.addField("1", []byte{})
	js.buf = append(js.buf, '1')
	if string(js.buf) != testString {
		t.Errorf("value in json is not valid. Expected: %s, received: %s.", testString, string(js.buf))
	}

	if js.buf[len(js.buf)-1] == '"' {
		t.Error("last byte in buf cannot be closing double quote")
	}

	testString = `"1": "1", "2": "2`
	js.buf = js.addField("2", js.buf)
	js.buf = append(js.buf, '2')
	if string(js.buf) != testString {
		t.Errorf("value in json is not valid. Expected: %s, received: %s.", testString, string(js.buf))
	}
}
