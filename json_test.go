package logg

import (
	"testing"
)

func TestLogg_getJson(t *testing.T) {
	js := newJSON()
	if js == nil {
		t.Error("js object is nil")
		t.FailNow()
	}

	if len(js.buf) == 0 {
		t.Error("js buffer is empty, but must contain a opening bracket '{'")
	}
}

func TestLogg_json_put(t *testing.T) {
	js := &json{
		buf: []byte{'0'},
	}

	js.put()

	js = jsonPool.Get().(*json)
	if len(js.buf) != 0 {
		t.Error("js buf must be empty after put")
	}
}

func TestLogg_json_close(t *testing.T) {
	js := &json{}
	js.close()

	if js.buf[len(js.buf)-1] != '}' {
		t.Error("js buf must end with closing bracket '}'")
	}
}

func TestLogg_json_addField(t *testing.T) {
	js := &json{}

	beforeLength := len(js.buf)
	js.addField("", []byte{})
	if beforeLength != len(js.buf) {
		t.Error("field with empty key must be missed")
	}

	testString := `"1": "1"`
	js.addField("1", []byte{'1'})
	if string(js.buf) != testString {
		t.Errorf("value in json is not valid. Expected: %s, received: %s.", testString, string(js.buf))
	}

	if js.buf[len(js.buf)-1] != '"' {
		t.Error("last byte in buf must be closing double quote")
	}

	testString = `"1": "1", "2": "2"`
	js.addField("2", []byte{'2'})
	if string(js.buf) != testString {
		t.Errorf("value in json is not valid. Expected: %s, received: %s.", testString, string(js.buf))
	}
}
