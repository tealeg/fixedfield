package fixedfield

import (
	"encoding/binary"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type readSpec struct {
	FieldValue reflect.Value
	FieldType reflect.StructField
	Length int
	Repeat int
	Encoding string
}

func (spec *readSpec) String() string {
	return fmt.Sprintf("Field Name: %s,\t Field Value: %v,\t Field Length: %d\n, repeat %d\n",
		spec.FieldType.Name, spec.FieldValue.Interface(), spec.Length, spec.Repeat)
}

func buildReadSpecs(structure interface{}) (readSpecs []readSpec, err error){
	var values, value reflect.Value 
	var spec readSpec
	var tag reflect.StructTag
	var length, repeat, encoding string

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
		encoding = tag.Get("encoding")
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
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if len(encoding) == 0 {
				spec.Encoding = "LE"
			} else {
				spec.Encoding = encoding
			}
		}
		readSpecs[i] = spec
	}
	return readSpecs, nil
}

func populateStructFromReadSpecAndBytes(target interface{}, readSpecs []readSpec, data io.Reader) (err error) {
	for _, spec := range readSpecs {
		var bytesRead int
		block := make([]byte, spec.Length)
		bytesRead, err = data.Read(block)
		if err != nil {
			return err
		}
		if bytesRead != spec.Length {
			return fmt.Errorf("Buffer underrun, %d of %d bytes read.", bytesRead, spec.Length)
		}
		switch spec.FieldValue.Kind() {
		case reflect.String:
			spec.FieldValue.SetString(string(block))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var intVal int
			var value int64
			switch strings.ToLower(spec.Encoding) {
			case "ascii":
				intVal, err = strconv.Atoi(string(block))
				if err != nil {
					return err
				}
				value = int64(intVal)
			case "bigendian":
				buffer := bytes.NewBuffer(block)
				binary.Read(buffer, binary.BigEndian, value)
			}
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
