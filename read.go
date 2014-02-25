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
	StructName  string
	Value       reflect.Value
	StructField reflect.StructField
	Length      int
	Repeat      int
	Encoding    string
	TrueBytes   []byte
	Children    []readSpec
}

// Return a string representation of the readSpec
func (spec *readSpec) String() string {
	return fmt.Sprintf(
		"Field Name: %s,\n"+
			"Field Value: %v\n"+
			"Field Length: %d\n"+
			"Repeat %d\n"+
			"Encoding %s\n"+
			"TrueBytes %s\n"+
			"Children %v\n",
		spec.StructField.Name, spec.Value.Interface(), spec.Length, spec.Repeat,
		spec.Encoding, string(spec.TrueBytes), spec.Children)
}

func getFieldLength(tag reflect.StructTag) (int, error) {
	var tagLength string
	tagLength = tag.Get("length")
	if len(tagLength) == 0 {
		return 1, nil
	}
	return strconv.Atoi(tagLength)
}

func getFieldRepeat(tag reflect.StructTag) (int, error) {
	var repeat string

	repeat = tag.Get("repeat")
	if len(repeat) == 0 {
		return 1, nil
	}
	return strconv.Atoi(repeat)
}

func getFieldEncoding(tag reflect.StructTag) string {
	var encoding string

	encoding = tag.Get("encoding")
	if len(encoding) == 0 {
		return "LE"
	}
	return encoding
}

func getFieldTrueBytes(tag reflect.StructTag) []byte {
	var trueChars string

	trueChars = tag.Get("trueChars")

	if len(trueChars) == 0 {
		return []byte("Yy")
	}
	return []byte(trueChars)
}

func buildReadSpecFromField(value reflect.Value, field reflect.StructField, structName string) (spec readSpec, err error) {
	var tag reflect.StructTag

	spec = readSpec{}
	spec.StructName = structName
	spec.Value = value
	spec.StructField = field
	tag = spec.StructField.Tag

	spec.Length, err = getFieldLength(tag)
	if err != nil {
		return spec, err
	}

	spec.Repeat, err = getFieldRepeat(tag)
	if err != nil {
		return spec, err
	}

	spec.Encoding = getFieldEncoding(tag)
	spec.TrueBytes = getFieldTrueBytes(tag)
	return spec, err
}

func buildReadSpecsFromStructValue(value reflect.Value, structName string) (readSpecs []readSpec, err error) {
	var fieldCount int
	var spec readSpec
	var subStructName string

	fieldCount = value.NumField()
	readSpecs = make([]readSpec, fieldCount)

	for i := 0; i < fieldCount; i++ {
		spec, err = buildReadSpecFromField(value.Field(i), value.Type().Field(i), structName)
		if err != nil {
			return nil, err
		}
		if spec.Value.Kind() == reflect.Struct {
			spec.Length = 0
			spec.Repeat = 0
			subStructName = spec.Value.Type().String()
			spec.Children, err = buildReadSpecsFromStructValue(
				spec.Value, subStructName)
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
	readSpecs, err = buildReadSpecsFromStructValue(value, structName)
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
	reflectType := spec.StructField.Type
	kind := reflectType.Kind()
	typeName := kind.String()
	name := spec.StructName + "." + spec.StructField.Name
	return fmt.Errorf("Failure unmarshalling %s field '%s'. Integer fields must be annotated with an encoding type of BigEndian, LittleEndian or ASCII", typeName, name)
}

// Given a readSpec, a block of bytes, and a block length, populate
// the field defined by the readSpec with a 64bit integer value
// encoded in the block of bytes.
func readInteger(spec readSpec, block []byte) (err error) {
	var value int64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		value, err = readASCIIInteger(block)
	case "bigendian", "be":
		value, err = readBinaryInteger(block, spec.Length, binary.BigEndian)
	case "littleendian", "le":
		value, err = readBinaryInteger(block, spec.Length, binary.LittleEndian)
	default:
		err = makeUnmarshalIntegerError(spec)
	}
	if err == nil {
		spec.Value.SetInt(value)
		return nil
	}
	return err
}

// Given a readSpec, a block of bytes, and a block length, populate
// the field defined by the readSpec with a 64bit unsigned integer value
// encoded in the block of bytes.
func readUnsignedInteger(spec readSpec, block []byte) (err error) {
	var intVal int
	var value uint64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		intVal, err = strconv.Atoi(string(block))
		if err == nil {
			value = uint64(intVal)
		}
	case "bigendian", "be":
		value, err = readBinaryUnsignedInteger(block, spec.Length, binary.BigEndian)
	case "littleendian", "le":
		value, err = readBinaryUnsignedInteger(block, spec.Length, binary.LittleEndian)
	default:
		err = makeUnmarshalIntegerError(spec)
	}
	if err == nil {
		spec.Value.SetUint(value)
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
func readFloat(spec readSpec, block []byte, kind reflect.Kind) (err error) {
	var f64Val float64
	switch strings.ToLower(spec.Encoding) {
	case "ascii":
		if kind == reflect.Float32 {
			f64Val, err = strconv.ParseFloat(string(block), 32)
		} else {
			f64Val, err = strconv.ParseFloat(string(block), 64)
		}
	case "bigendian", "be":
		f64Val, err = readBinaryFloat(block, spec.Length, binary.BigEndian)
	case "littleendian", "le":
		f64Val, err = readBinaryFloat(block, spec.Length, binary.LittleEndian)
	default:
		err = fmt.Errorf("Invalid encoding for a floating point value specified. %s",
			spec.String())
	}

	if err == nil {
		spec.Value.SetFloat(f64Val)
	}
	return err
}

// Read a boolean from a block of bytes using encoding
// information from the readSpec.
func readBool(spec readSpec, block []byte) (err error) {
	var boolVal bool

	switch strings.ToLower(spec.Encoding) {
	case "littleendian", "le", "bigendian", "be", "byte":
		if spec.Length > 1 {
			err = fmt.Errorf("Booleans can only be 1 byte long, %d bytes specified for %s", spec.Length, spec.StructField.Name)
		}
		boolVal = int(block[0]) != 0
	case "ascii":
		boolVal = bytes.Contains(spec.TrueBytes, block)
	default:
		err = fmt.Errorf("Invalid encoding for a boolean value specified. %s",
			spec.String())
	}
	if err == nil {
		spec.Value.SetBool(boolVal)
	}
	return err

}

func populateKind(kind reflect.Kind, block []byte, spec readSpec, data io.Reader) (err error) {
	switch kind {
	case reflect.String:
		spec.Value.SetString(string(block))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = readInteger(spec, block)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = readUnsignedInteger(spec, block)
	case reflect.Float64, reflect.Float32:
		err = readFloat(spec, block, kind)
	case reflect.Bool:
		err = readBool(spec, block)
	case reflect.Struct:
		// Recur, exploring the nested specification.
		err = populateStructFromReadSpecAndBytes(spec.Children, data)
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

// Given a slice of readSpecs and some data, populate the target
// struct elements from the data.
func populateStructFromReadSpecAndBytes(readSpecs []readSpec, data io.Reader) (err error) {
	var block []byte
	var sliceType reflect.Type
	var elemKind reflect.Kind

	for _, spec := range readSpecs {
		kind := spec.Value.Kind()
		if kind == reflect.Slice {
			sliceType = spec.Value.Type()
			elemKind = sliceType.Elem().Kind()
			if !spec.Value.CanSet() {
				return fmt.Errorf("Cannot set slice, %s", spec.StructName)
			}
			spec.Value.Set(
				reflect.MakeSlice(sliceType, spec.Repeat, spec.Repeat))
			sliceValue := spec.Value
			for offset := 0; offset < spec.Repeat; offset++ {
				block, err = readBlock(data, spec.Length)
				if err != nil {
					return err
				}
				spec.Value = sliceValue.Index(offset)
				err = populateKind(elemKind, block, spec, data)
			}
		} else {
			block, err = readBlock(data, spec.Length)
			if err != nil {
				return err
			}
			err = populateKind(kind, block, spec, data)
		}

		if err != nil {
			return err
		}
		// Struct

		// Array
	}
	return nil
}

func Unmarshal(data []byte, v interface{}) (err error) {
	readSpec, err := buildReadSpecs(v)
	if err != nil {
		return err
	}
	return populateStructFromReadSpecAndBytes(readSpec, bytes.NewBuffer(data))

}
