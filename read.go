package fixedfield

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

type readSpec struct {
	FieldValue reflect.Value
	FieldType  reflect.StructField
	Length     int
	Repeat     int
	Encoding   string
}

func (spec *readSpec) String() string {
	return fmt.Sprintf("Field Name: %s,\t Field Value: %v,\t Field Length: %d\n, repeat %d\n",
		spec.FieldType.Name, spec.FieldValue.Interface(), spec.Length, spec.Repeat)
}

func buildReadSpecs(structure interface{}) (readSpecs []readSpec, err error) {
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
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if len(encoding) == 0 {
				spec.Encoding = "LE"
			} else {
				spec.Encoding = encoding
			}
		case reflect.Float64, reflect.Float32:
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

func readBinaryInteger(block []byte, blockLength int, byteOrder binary.ByteOrder) (value int64, err error) {
	buffer := bytes.NewBuffer(block)
	switch blockLength {
	case 1:
		var val int8
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = int64(val)
		}
	case 2:
		var val int16
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = int64(val)
		}
	case 4:
		var val int32
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = int64(val)
		}
	case 8:
		err = binary.Read(buffer, byteOrder, &value)
	}
	return value, err
}

func readBinaryUnsignedInteger(block []byte, blockLength int, byteOrder binary.ByteOrder) (value uint64, err error) {
	buffer := bytes.NewBuffer(block)
	switch blockLength {
	case 1:
		var val uint8
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = uint64(val)
		}
	case 2:
		var val uint16
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = uint64(val)
		}
	case 4:
		var val uint32
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = uint64(val)
		}
	case 8:
		err = binary.Read(buffer, byteOrder, &value)
	}
	return value, err
}

func readASCIIInteger(block []byte) (value int64, err error) {
	var intVal int
	intVal, err = strconv.Atoi(string(block))
	if err != nil {
		return
	}
	value = int64(intVal)
	return
}

func readInteger(spec readSpec, block []byte, blockLength int) (err error) {
	var value int64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		value, err = readASCIIInteger(block)
	case "bigendian", "be":
		value, err = readBinaryInteger(block, blockLength, binary.BigEndian)
	case "litteendian", "le":
		value, err = readBinaryInteger(block, blockLength, binary.LittleEndian)
	}
	if err == nil {
		spec.FieldValue.SetInt(value)
		return nil
	}
	return err
}

func readUnsignedInteger(spec readSpec, block []byte, blockLength int) (err error) {
	var intVal int
	var value uint64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		intVal, err = strconv.Atoi(string(block))
		if err == nil {
			value = uint64(intVal)
		}
	case "bigendian", "be":
		value, err = readBinaryUnsignedInteger(block, blockLength, binary.BigEndian)
	case "litteendian", "le":
		value, err = readBinaryUnsignedInteger(block, blockLength, binary.LittleEndian)
	}
	if err == nil {
		spec.FieldValue.SetUint(value)
		return nil
	}
	return err
}

func readFloat(spec readSpec, block []byte, bytesRead int, kind reflect.Kind) (err error) {
	var f64Val float64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		if kind == reflect.Float32 {
			f64Val, err = strconv.ParseFloat(string(block), 32)
		} else {
			f64Val, err = strconv.ParseFloat(string(block), 64)
		}
	}
	if err == nil {
		spec.FieldValue.SetFloat(f64Val)
	}
	return err
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
		kind := spec.FieldValue.Kind()
		switch kind {
		case reflect.String:
			spec.FieldValue.SetString(string(block))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = readInteger(spec, block, bytesRead)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = readUnsignedInteger(spec, block, bytesRead)
		case reflect.Float64, reflect.Float32:
			err = readFloat(spec, block, bytesRead, kind)
		}

		if err != nil {
			return err
		}
		// Invalid Kind = iota
		// Bool
		// Uintptr
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
