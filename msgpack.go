package msgpack

import (
	"bytes"
	"math"
	"reflect"
)

const (
	maxPositiveFixnum = 0x7f
	maxFixstrLength   = 31
	maxStr8Length     = 1<<8 - 1
	maxStr16Length    = 1<<16 - 1
)

func Marshal(i interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if i == nil {
		buf.WriteByte(0xc0)
	} else if err := pack(reflect.ValueOf(i), &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func pack(v reflect.Value, buf *bytes.Buffer) error {
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
		if ival >= 0 {
			writeUint(uint64(ival), buf)
		} else {
			writeInt(ival, buf)
		}
	case reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		writeUint(v.Uint(), buf)
	case reflect.Float32: // TODO: treat integral floats as integers
		fval := v.Float()
		floatBits := math.Float32bits(float32(fval))
		writeUint32(0xca, floatBits, buf)
	case reflect.Float64:
		fval := v.Float()
		floatBits := math.Float64bits(fval)
		writeUint64(0xcb, floatBits, buf)
	case reflect.String:
		str := v.String()
		length := len(str)
		writeStringHeader(length, buf)
		buf.WriteString(str)
	}

	return nil
}

func writeStringHeader(length int, buf *bytes.Buffer) {
	if length <= maxFixstrLength {
		buf.WriteByte(byte(160 + length)) // 160 == 0xA0 == 1010 0000
	} else if length <= maxStr8Length {
		buf.WriteByte(0xd9)
		buf.WriteByte(byte(length))
	} else if length <= maxStr16Length {
		buf.WriteByte(0xda)
		i16 := int16(length)
		buf.WriteByte(byte(i16 >> 8))
		buf.WriteByte(byte(i16 & 0x00FF))
	}
}

func writeUint32(prefix byte, uval uint32, buf *bytes.Buffer) {
	buf.WriteByte(prefix)
	buf.WriteByte(byte(uval >> 24))
	buf.WriteByte(byte((uval & 0x00FF0000) >> 16))
	buf.WriteByte(byte((uval & 0x0000FF00) >> 8))
	buf.WriteByte(byte((uval & 0x000000FF)))
}

func writeUint64(prefix byte, uval uint64, buf *bytes.Buffer) {
	buf.WriteByte(prefix)
	buf.WriteByte(byte(uval >> 56))
	buf.WriteByte(byte((uval & 0x00FF000000000000) >> 48))
	buf.WriteByte(byte((uval & 0x0000FF0000000000) >> 40))
	buf.WriteByte(byte((uval & 0x000000FF00000000) >> 32))
	buf.WriteByte(byte((uval & 0x00000000FF000000) >> 24))
	buf.WriteByte(byte((uval & 0x0000000000FF0000) >> 16))
	buf.WriteByte(byte((uval & 0x000000000000FF00) >> 8))
	buf.WriteByte(byte((uval & 0x00000000000000FF)))
}

func writeUint(uval uint64, buf *bytes.Buffer) {
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
		writeUint32(0xce, u32, buf)
	} else {
		writeUint64(0xcf, uval, buf)
	}
}

func writeInt(ival int64, buf *bytes.Buffer) {
	if isFixnum(ival) {
		buf.WriteByte(byte(ival))
		return
	} else if ival > math.MinInt8 && ival <= math.MaxInt8 {
		buf.WriteByte(0xd0)
		buf.WriteByte(byte(ival))
	} else if ival > math.MinInt16 && ival <= math.MaxInt16 {
		i16 := int16(ival)
		buf.WriteByte(0xd1)
		buf.WriteByte(byte(i16 >> 8))
		buf.WriteByte(byte(i16 & 0x00FF))
	} else if ival > math.MinInt32 && ival <= math.MaxInt32 {
		i32 := int32(ival)
		buf.WriteByte(0xd2)
		buf.WriteByte(byte(i32 >> 24))
		buf.WriteByte(byte((i32 & 0x00FF0000) >> 16))
		buf.WriteByte(byte((i32 & 0x0000FF00) >> 8))
		buf.WriteByte(byte((i32 & 0x000000FF)))
	} else {
		buf.WriteByte(0xd3)
		buf.WriteByte(byte(ival >> 56))
		buf.WriteByte(byte((ival & 0x00FF000000000000) >> 48))
		buf.WriteByte(byte((ival & 0x0000FF0000000000) >> 40))
		buf.WriteByte(byte((ival & 0x000000FF00000000) >> 32))
		buf.WriteByte(byte((ival & 0x00000000FF000000) >> 24))
		buf.WriteByte(byte((ival & 0x0000000000FF0000) >> 16))
		buf.WriteByte(byte((ival & 0x000000000000FF00) >> 8))
		buf.WriteByte(byte((ival & 0x00000000000000FF)))
	}
}

func isFixnum(n int64) bool {
	return (n >= 0 && n <= maxPositiveFixnum) || (n < 0 && n >= -32)
}
