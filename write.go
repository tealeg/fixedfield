package fixedfield

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func marshalASCIIInteger(s spec) (block []byte, err error) {
	var formatString, candidate string

	formatString = "%" + strconv.Itoa(s.Length) + "d"
	candidate = fmt.Sprintf(formatString, int(s.Value.Int()))
	if len(candidate) > s.Length {
		return nil, fmt.Errorf("Field %s.%s overflowed configured field length (Tried to write %s to a %d length ASCII field)",
			s.StructName, s.StructField.Name, candidate, s.Length)
	}
	return []byte(candidate), nil
}

func marshalInteger(s spec) (block []byte, err error) {
	switch strings.ToLower(s.Encoding) {
	case "ascii":
		return marshalASCIIInteger(s)
	}
	return
}

func marshalKind(kind reflect.Kind, s spec) (block []byte, err error) {
	switch kind {
	case reflect.String:
		block = []byte(s.Value.String())
	case reflect.Int:
		block, err = marshalInteger(s)
	}
	return block, err
}


func populateBytesFromSpecAndStruct(specs []spec) (data []byte, err error) {
	var buffer *bytes.Buffer
	var block []byte

	buffer = bytes.NewBuffer(nil)
	for _, s := range specs {
		kind := s.Value.Kind()
		block, err = marshalKind(kind, s)
		if err != nil {
			return nil, err
		}
		_, err = buffer.Write(block)
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}


func Marshal(v interface{}) (result []byte, err error) {
	var specs []spec

	specs, err = buildSpecs(v)
	if err != nil {
		return nil, err
	}
	return populateBytesFromSpecAndStruct(specs)
}
