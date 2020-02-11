package logg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
	"unsafe"
)

func write(w io.Writer, b []byte) error {
	n, err := w.Write(b)
	if err != nil {
		return err
	}

	if n != len(b) {
		return io.ErrShortWrite
	}

	return nil
}

func caller(calldepth int, shortFile bool) (file string, line int) {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	if shortFile {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	return
}

func timestamp(t time.Time, flags int) (buf []byte) {
	if t.IsZero() {
		buf = append(buf, "0000-00-00 00:00:00"...)
		return
	}

	t = t.Add(time.Nanosecond * 500) // To round under microsecond
	year := t.Year()
	year100 := year / 100
	year1 := year % 100
	month := t.Month()
	day := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	micro := t.Nanosecond() / 1000

	if flags&log.Ldate != 0 && flags&log.Ltime != 0 {
		buf = append(buf, []byte{
			digits10[year100], digits01[year100],
			digits10[year1], digits01[year1],
			'-',
			digits10[month], digits01[month],
			'-',
			digits10[day], digits01[day],
			' ',
			digits10[hour], digits01[hour],
			':',
			digits10[minute], digits01[minute],
			':',
			digits10[second], digits01[second],
		}...)
	} else if flags&log.Ldate != 0 && flags&log.Ltime == 0 {
		buf = append(buf, []byte{
			digits10[year100], digits01[year100],
			digits10[year1], digits01[year1],
			'-',
			digits10[month], digits01[month],
			'-',
			digits10[day], digits01[day],
		}...)
	} else if flags&log.Ldate == 0 && flags&log.Ltime != 0 {
		buf = append(buf, []byte{
			digits10[hour], digits01[hour],
			':',
			digits10[minute], digits01[minute],
			':',
			digits10[second], digits01[second],
		}...)
	}

	if flags&log.Lmicroseconds != 0 {
		micro10000 := micro / 10000
		micro100 := micro / 100 % 100
		micro1 := micro % 100
		buf = append(buf, []byte{
			'.',
			digits10[micro10000], digits01[micro10000],
			digits10[micro100], digits01[micro100],
			digits10[micro1], digits01[micro1],
		}...)
	}

	return
}

func defineLevel(data *[]byte) (lvl level) {
	lvl = Empty

	if len(*data) == 0 {
		return
	}

	if (*data)[0] == '[' {
		y := bytes.IndexByte(*data, ']')
		if y > 0 {
			*data = (*data)[1:]
			*data = append((*data)[:y-1], (*data)[y:]...)
		}
	}

	if (*data)[0] == 'D' || (*data)[0] == 'I' || (*data)[0] == 'E' || (*data)[0] == 'W' || (*data)[0] == 'P' {
		for i := 0; i < len(longLevels); i++ {
			l := longLevels[i]
			if bytes.Contains((*data)[:5], []byte(l)) {
				lvl = level(i)
				x := bytes.Index(*data, []byte(l))
				*data = (*data)[x+len([]byte(l)):]
				break
			}
		}

		if lvl == Empty {
			for i := 0; i < len(levels); i++ {
				l := levels[i]
				if bytes.Contains((*data)[:3], []byte(l)) {
					lvl = level(i)
					x := bytes.Index(*data, []byte(l))
					*data = (*data)[x+len([]byte(l)):]
					break
				}
			}
		}
	}

	for (*data)[0] == 32 { // trip space
		*data = append((*data)[:0], (*data)[1:]...)
	}

	return
}

func colorize(buf *[]byte, color int) {
	*buf = append([]byte(fmt.Sprintf("%v[%dm", escape, 30+color)), *buf...)
	*buf = append(*buf, []byte(escapeClose)...)
}

func unsafeCompare(a string, b []byte) int {
	abp := *(*[]byte)(unsafe.Pointer(&a))
	return bytes.Compare(abp, b)
}

func unsafeEqual(a string, b []byte) bool {
	bbp := *(*string)(unsafe.Pointer(&b))
	return a == bbp
}
