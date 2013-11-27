package fixedfield

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type readSpec struct {
	FieldValue reflect.Value
	FieldType reflect.StructField
	Length int
	Repeat int
}

func (spec *readSpec) String() string {
	return fmt.Sprintf("Field Name: %s,\t Field Value: %v,\t Field Length: %d\n, repeat %d\n",
		spec.FieldType.Name, spec.FieldValue.Interface(), spec.Length, spec.Repeat)
}

func buildReadSpecs(structure interface{}) (readSpecs []readSpec, err error){
	var values, value reflect.Value 
	var spec readSpec
	var tag reflect.StructTag
	var length, repeat string

	values = reflect.ValueOf(structure).Elem()
	readSpecs = make([]readSpec, values.NumField())

	for i := 0; i < values.NumField(); i++ {
		spec = readSpec{}
		value = values.Field(i)
		spec.FieldValue = value
		spec.FieldType = values.Type().Field(i)
		tag = spec.FieldType.Tag
		length = tag.Get("length")
		repeat = tag.Get("repeat")
		if len(length) == 0 {
			spec.Length = 0
		} else {
			spec.Length, err = strconv.Atoi(length)
			if err != nil {
				return nil, err
			}
		}
		if len(repeat) == 0 {
			spec.Repeat = 1
		} else {
			spec.Repeat, err = strconv.Atoi(repeat)
			if err != nil {
				return nil, err
			}
		}
		readSpecs[i] = spec
	}
	return readSpecs, nil
}

func populateStructFromReadSpecAndBytes(target interface{}, readSpecs []readSpec, data io.Reader) error {
	for _, spec := range readSpecs {
		block := make([]byte, spec.Length)
		n, err := data.Read(block)
		if err != nil {
			return err
		}
		if n != spec.Length {
			return fmt.Errorf("Buffer underrun, %d of %d bytes read.", n, spec.Length)
		}
		switch spec.FieldValue.Kind() {
		case reflect.String:
			spec.FieldValue.SetString(string(block))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var value int64
			buffer := bytes.NewBuffer(block)
			binary.Read(buffer, binary.BigEndian, value)
			spec.FieldValue.SetInt(value)
		}
        // Invalid Kind = iota
        // Bool
        // Int
        // Int8
        // Int16
        // Int32
        // Int64
        // Uint
        // Uint8
        // Uint16
        // Uint32
        // Uint64
        // Uintptr
        // Float32
        // Float64
        // Complex64
        // Complex128
        // Array
        // Chan
        // Func
        // Interface
        // Map
        // Ptr
        // Slice

        // Struct
        // UnsafePointer
	}
	return nil
}
