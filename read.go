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

// A readSpec is created, by buildReadSpecs, for each field in a
// target structure we wish to populated.  These readSpecs are used by
// populateStructFromReadSpecAndByte to guide the unmarshalling of
// byte data into the target struct.
type readSpec struct {
	StructName string
	FieldValue reflect.Value
	FieldType  reflect.StructField
	Length     int
	Repeat     int
	Encoding   string
	TrueBytes  []byte
	Children []readSpec
}

func (spec *readSpec) String() string {
	return fmt.Sprintf("Field Name: %s,\t Field Value: %v,\t Field Length: %d\n, repeat %d\n",
		spec.FieldType.Name, spec.FieldValue.Interface(), spec.Length, spec.Repeat)
}


func buildReadSpecsFromElems(value reflect.Value, structName string) (readSpecs []readSpec, err error){
	var fieldCount int
	var spec readSpec
	var tag reflect.StructTag
	var length, repeat, encoding, trueChars string
	var subStructName string

	fieldCount = value.NumField()
	readSpecs = make([]readSpec, fieldCount)

	for i := 0; i < fieldCount; i++ {
		spec = readSpec{}
		spec.StructName = structName
		spec.FieldValue = value.Field(i)
		spec.FieldType = value.Type().Field(i)
		tag = spec.FieldType.Tag
		length = tag.Get("length")
		repeat = tag.Get("repeat")
		encoding = tag.Get("encoding")
		trueChars = tag.Get("trueChars")
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
		switch spec.FieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64:
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
		case reflect.Bool:
			if spec.Length == 0 {
				spec.Length = 1
			}
			if len(encoding) == 0 {
				spec.Encoding = "LE"
			} else {
				spec.Encoding = encoding
				if encoding == "ascii" {
					if len(trueChars) == 0 {
						spec.TrueBytes = []byte("Yy")
					} else {
						spec.TrueBytes = []byte(trueChars)
					}
				}
			}
		case reflect.Struct:
			subStructName = spec.FieldValue.Type().String()
			spec.Children, err = buildReadSpecsFromElems(
				spec.FieldValue, subStructName)
			if err != nil {
				return nil, err
			}
		}
		readSpecs[i] = spec
	}
	return readSpecs, nil
}

// Convert annotation on a structure into a specification for what
// should be read from a fixed field file.
func buildReadSpecs(structure interface{}) (readSpecs []readSpec, err error) {
	var structValue, value reflect.Value
	var structType reflect.Type
	var structName string

	structValue = reflect.ValueOf(structure)
	structType = reflect.TypeOf(structure)
	structName = structType.String()

	value = structValue.Elem()
	readSpecs, err = buildReadSpecsFromElems(value, structName)
	return readSpecs, nil
}

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
	intVal, err = strconv.Atoi(blockString)
	intVal = intVal * multiple
	if err != nil {
		return
	}
	value = int64(intVal)
	return
}

func makeUnmarshalIntegerError(spec readSpec) error {
	reflectType := spec.FieldType.Type
	kind := reflectType.Kind()
	typeName := kind.String()
	name := spec.StructName + "." + spec.FieldType.Name
	return fmt.Errorf("Failure unmarshalling %s field '%s'. Integer fields must be annotated with an encoding type of BigEndian, LittleEndian or ASCII", typeName, name)
}

// Given a readSpec, a block of bytes, and a block length, populate
// the field defined by the readSpec with a 64bit integer value
// encoded in the block of bytes.
func readInteger(spec readSpec, block []byte, blockLength int) (err error) {
	var value int64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		value, err = readASCIIInteger(block)
	case "bigendian", "be":
		value, err = readBinaryInteger(block, blockLength, binary.BigEndian)
	case "littleendian", "le":
		value, err = readBinaryInteger(block, blockLength, binary.LittleEndian)
	default:
		err = makeUnmarshalIntegerError(spec)
	}
	if err == nil {
		spec.FieldValue.SetInt(value)
		return nil
	}
	return err
}

// Given a readSpec, a block of bytes, and a block length, populate
// the field defined by the readSpec with a 64bit unsigned integer value
// encoded in the block of bytes.
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
	case "littleendian", "le":
		value, err = readBinaryUnsignedInteger(block, blockLength, binary.LittleEndian)
	default:
		err = makeUnmarshalIntegerError(spec)
	}
	if err == nil {
		spec.FieldValue.SetUint(value)
		return nil
	}
	return err
}

// Given a block of bytes, a block length, and a byte order, populate
// the field defined by the readSpec with a 64bit float value encoded
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
// information from the readSpec.
func readFloat(spec readSpec, block []byte, bytesRead int, kind reflect.Kind) (err error) {
	var f64Val float64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		if kind == reflect.Float32 {
			f64Val, err = strconv.ParseFloat(string(block), 32)
		} else {
			f64Val, err = strconv.ParseFloat(string(block), 64)
		}
	case "bigendian", "be":
		f64Val, err = readBinaryFloat(block, bytesRead, binary.BigEndian)
	case "littleendian", "le":
		f64Val, err = readBinaryFloat(block, bytesRead, binary.LittleEndian)
	default:
		err = fmt.Errorf("Invalid encoding for a floating point value specified. %s",
			spec.String())
	}

	if err == nil {
		spec.FieldValue.SetFloat(f64Val)
	}
	return err
}

// Read a boolean from a block of bytes using encoding
// information from the readSpec.
func readBool(spec readSpec, block []byte, bytesRead int) (err error) {
	var boolVal bool

	switch strings.ToLower(spec.Encoding) {
	case "littleendian", "le", "bigendian", "be", "byte":
		if bytesRead > 1 {
			err = fmt.Errorf("Booleans can only be 1 byte long, %d bytes specified for %s", spec.Length, spec.FieldType.Name)
		}
		boolVal = int(block[0]) != 0
	case "ascii":
		boolVal = bytes.Contains(spec.TrueBytes, block)
	default:
		err = fmt.Errorf("Invalid encoding for a boolean value specified. %s",
			spec.String())
	}
	if err == nil {
		spec.FieldValue.SetBool(boolVal)
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
		case reflect.Bool:
			err = readBool(spec, block, bytesRead)
		}

		if err != nil {
			return err
		}
		// Struct

		// Array
	}
	return nil
}
