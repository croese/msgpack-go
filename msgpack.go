package msgpack

import (
	"bytes"
	"math"
	"reflect"
)

const maxPositiveFixnum = 0x7f

func Marshal(i interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if i == nil {
		buf.WriteByte(0xc0)
	} else {
		v := reflect.ValueOf(i)
		switch v.Kind() {
		case reflect.Bool:
			if v.Bool() {
				buf.WriteByte(0xc3)
			} else {
				buf.WriteByte(0xc2)
			}
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			ival := v.Int()
			if isFixnum(ival) {
				buf.WriteByte(byte(ival))
			} else if ival >= 0 {
				WriteUint(uint64(ival), &buf)
			} else {
				buf.WriteByte(0xd0)
				buf.WriteByte(byte(ival))
			}
		case reflect.Uint,
			reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			WriteUint(v.Uint(), &buf)
		}
	}

	return buf.Bytes(), nil
}

func WriteUint(uval uint64, buf *bytes.Buffer) {
	if uval <= maxPositiveFixnum {
		buf.WriteByte(byte(uval))
	} else if uval <= math.MaxUint8 {
		buf.WriteByte(0xcc)
		buf.WriteByte(byte(uval))
	} else if uval <= math.MaxUint16 {
		u16 := uint16(uval)
		buf.WriteByte(0xcd)
		buf.WriteByte(byte(u16 >> 8))
		buf.WriteByte(byte(u16 & 0x00FF))
	} else if uval <= math.MaxUint32 {
		u32 := uint32(uval)
		buf.WriteByte(0xce)
		buf.WriteByte(byte(u32 >> 24))
		buf.WriteByte(byte((u32 & 0x00FF0000) >> 16))
		buf.WriteByte(byte((u32 & 0x0000FF00) >> 8))
		buf.WriteByte(byte((u32 & 0x000000FF)))
	}
	buf.WriteByte(0xcf)
	buf.WriteByte(byte(uval >> 56))
	buf.WriteByte(byte((uval & 0x00FF000000000000) >> 48))
	buf.WriteByte(byte((uval & 0x0000FF0000000000) >> 40))
	buf.WriteByte(byte((uval & 0x000000FF00000000) >> 32))
	buf.WriteByte(byte((uval & 0x00000000FF000000) >> 24))
	buf.WriteByte(byte((uval & 0x0000000000FF0000) >> 16))
	buf.WriteByte(byte((uval & 0x000000000000FF00) >> 8))
	buf.WriteByte(byte((uval & 0x00000000000000FF)))
}

func isFixnum(n int64) bool {
	return (n >= 0 && n <= maxPositiveFixnum) || (n < 0 && n >= -32)
}
