package msgpack

import (
	"bytes"
	"fmt"
	"reflect"
)

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
			fmt.Printf("DEBUG: %d\n", v.Int())
		}
	}

	return buf.Bytes(), nil
}
