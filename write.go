package fixedfield

import (
	"bytes"
	"reflect"
)

func marshalKind(kind reflect.Kind, s spec) (block []byte, err error) {
	switch kind {
	case reflect.String:
		block = []byte(s.Value.String())
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
