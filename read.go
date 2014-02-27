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


// Convert an array of bytes, of a known length and byte order, into
// a signed 64 bit integer.
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

// Convert an array of bytes, of a known length and byte order, into
// an unsigned 64 bit integer.
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

// Convert an array of ASCII chars, of a known length, into a 64 bit integer.
func readASCIIInteger(block []byte) (value int64, err error) {
	var intVal int
	var blockString string = string(block)
	var multiple int = 1
	if blockString[0] == '-' {
		blockString = blockString[1:]
		multiple = -1
	}
	intVal, err = strconv.Atoi(strings.TrimSpace(blockString))
	intVal = intVal * multiple
	if err != nil {
		return
	}
	value = int64(intVal)
	return
}

func makeUnmarshalIntegerError(s spec) error {
	reflectType := s.StructField.Type
	kind := reflectType.Kind()
	typeName := kind.String()
	name := s.StructName + "." + s.StructField.Name
	return fmt.Errorf("Failure unmarshalling %s field '%s'. Integer fields must be annotated with an encoding type of BigEndian, LittleEndian or ASCII", typeName, name)
}

// Given a spec, a block of bytes, and a block length, populate
// the field defined by the spec with a 64bit integer value
// encoded in the block of bytes.
func readInteger(s spec, block []byte) (err error) {
	var value int64
	switch strings.ToLower(s.Encoding) {
	case "ascii":
		value, err = readASCIIInteger(block)
	case "bigendian", "be":
		value, err = readBinaryInteger(block, s.Length, binary.BigEndian)
	case "littleendian", "le":
		value, err = readBinaryInteger(block, s.Length, binary.LittleEndian)
	default:
		err = makeUnmarshalIntegerError(s)
	}
	if err == nil {
		s.Value.SetInt(value)
		return nil
	}
	return err
}

// Given a spec, a block of bytes, and a block length, populate
// the field defined by the spec with a 64bit unsigned integer value
// encoded in the block of bytes.
func readUnsignedInteger(s spec, block []byte) (err error) {
	var intVal int
	var value uint64
	switch strings.ToLower(s.Encoding) {
	case "ascii":
		intVal, err = strconv.Atoi(string(block))
		if err == nil {
			value = uint64(intVal)
		}
	case "bigendian", "be":
		value, err = readBinaryUnsignedInteger(block, s.Length, binary.BigEndian)
	case "littleendian", "le":
		value, err = readBinaryUnsignedInteger(block, s.Length, binary.LittleEndian)
	default:
		err = makeUnmarshalIntegerError(s)
	}
	if err == nil {
		s.Value.SetUint(value)
		return nil
	}
	return err
}

// Given a block of bytes, a block length, and a byte order, populate
// the field defined by the spec with a 64bit float value encoded
// in the block of bytes.
func readBinaryFloat(block []byte, blockLength int, byteOrder binary.ByteOrder) (value float64, err error) {
	buffer := bytes.NewBuffer(block)
	switch blockLength {
	case 4:
		var val float32
		err = binary.Read(buffer, byteOrder, &val)
		if err == nil {
			value = float64(val)
		}
	case 8:
		err = binary.Read(buffer, byteOrder, &value)
	default:
		err = fmt.Errorf("Binary floats must have a length of either 4 or 8 bytes (float32 or float64 respectively).")
	}
	return
}

// Read a 64bit float from a block of bytes using encoding
// information from the spec.
func readFloat(s spec, block []byte, kind reflect.Kind) (err error) {
	var f64Val float64
	switch strings.ToLower(s.Encoding) {
	case "ascii":
		if kind == reflect.Float32 {
			f64Val, err = strconv.ParseFloat(string(block), 32)
		} else {
			f64Val, err = strconv.ParseFloat(string(block), 64)
		}
	case "bigendian", "be":
		f64Val, err = readBinaryFloat(block, s.Length, binary.BigEndian)
	case "littleendian", "le":
		f64Val, err = readBinaryFloat(block, s.Length, binary.LittleEndian)
	default:
		err = fmt.Errorf("Invalid encoding for a floating point value specified. %s",
			s.String())
	}

	if err == nil {
		s.Value.SetFloat(f64Val)
	}
	return err
}

// Read a boolean from a block of bytes using encoding
// information from the spec.
func readBool(s spec, block []byte) (err error) {
	var boolVal bool

	switch strings.ToLower(s.Encoding) {
	case "littleendian", "le", "bigendian", "be", "byte":
		if s.Length > 1 {
			err = fmt.Errorf("Booleans can only be 1 byte long, %d bytes specified for %s", s.Length, s.StructField.Name)
		}
		boolVal = int(block[0]) != 0
	case "ascii":
		boolVal = bytes.Contains(s.TrueBytes, block)
	default:
		err = fmt.Errorf("Invalid encoding for a boolean value specified. %s",
			s.String())
	}
	if err == nil {
		s.Value.SetBool(boolVal)
	}
	return err
}

func populateKind(kind reflect.Kind, block []byte, s spec, data io.Reader) (err error) {
	switch kind {
	case reflect.String:
		s.Value.SetString(string(block))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = readInteger(s, block)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = readUnsignedInteger(s, block)
	case reflect.Float64, reflect.Float32:
		err = readFloat(s, block, kind)
	case reflect.Bool:
		err = readBool(s, block)
	case reflect.Struct:
		// Recur, exploring the nested specification.
		err = populateStructFromSpecAndBytes(s.Children, data)
	}
	return err
}

// Read an array of bytes of given lengdh from the provided data.
func readBlock(data io.Reader, length int) (block []byte, err error) {
	var bytesRead int
	block = make([]byte, length)
	bytesRead, err = data.Read(block)
	if bytesRead != length {
		return nil, fmt.Errorf("Buffer underrun, %d of %d bytes read.", bytesRead, length)
	}
	return block, err
}

// Given a slice of specs and some data, populate the target
// struct elements from the data.
func populateStructFromSpecAndBytes(specs []spec, data io.Reader) (err error) {
	var block []byte
	var sliceType reflect.Type
	var elemKind reflect.Kind

	for _, s := range specs {
		kind := s.Value.Kind()
		if kind == reflect.Slice {
			sliceType = s.Value.Type()
			elemKind = sliceType.Elem().Kind()
			if !s.Value.CanSet() {
				return fmt.Errorf("Cannot set slice, %s", s.StructName)
			}
			s.Value.Set(
				reflect.MakeSlice(sliceType, s.Repeat, s.Repeat))
			sliceValue := s.Value
			for offset := 0; offset < s.Repeat; offset++ {
				block, err = readBlock(data, s.Length)
				if err != nil {
					return err
				}
				s.Value = sliceValue.Index(offset)
				err = populateKind(elemKind, block, s, data)
			}
		} else {
			block, err = readBlock(data, s.Length)
			if err != nil {
				return err
			}
			err = populateKind(kind, block, s, data)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func Unmarshal(data []byte, v interface{}) (err error) {
	var specs []spec

	specs, err = buildSpecs(v)
	if err != nil {
		return err
	}
	return populateStructFromSpecAndBytes(specs, bytes.NewBuffer(data))

}
