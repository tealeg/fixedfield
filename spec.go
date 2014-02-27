package fixedfield

import (
	"fmt"
	"reflect"
	"strconv"
)

// A spec is created, by buildSpecs, for each field in a
// target structure we wish to populated.  These specs are used by
// populateStructFromSpecAndByte to guide the unmarshalling of
// byte data into the target struct.
type spec struct {
	StructName  string
	Value       reflect.Value
	StructField reflect.StructField
	Length      int
	Repeat      int
	Encoding    string
	Padding     string
	TrueBytes   []byte
	Children    []spec
}

// Return a string representation of the spec
func (s *spec) String() string {
	return fmt.Sprintf(
		"Field Name: %s,\n"+
			"Field Value: %v\n"+
			"Field Length: %d\n"+
			"Repeat %d\n"+
			"Encoding %s\n"+
			"TrueBytes %s\n"+
			"Children %v\n",
		s.StructField.Name, s.Value.Interface(), s.Length, s.Repeat,
		s.Encoding, string(s.TrueBytes), s.Children)
}

func (s *spec) Size() int {
	return s.Length * s.Repeat
}


func getPadding(tag reflect.StructTag) string {
	var padding string
	padding = tag.Get("padding")
	if len(padding) == 0 {
		padding = "0"
	}
	return padding
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

func buildSpecFromField(value reflect.Value, field reflect.StructField, structName string) (s spec, err error) {
	var tag reflect.StructTag

	s = spec{}
	s.StructName = structName
	s.Value = value
	s.StructField = field
	tag = s.StructField.Tag

	s.Length, err = getFieldLength(tag)
	if err != nil {
		return s, err
	}

	s.Repeat, err = getFieldRepeat(tag)
	if err != nil {
		return s, err
	}

	s.Encoding = getFieldEncoding(tag)
	s.TrueBytes = getFieldTrueBytes(tag)
	s.Padding = getPadding(tag)
	return s, err
}

func buildSpecsFromStructValue(value reflect.Value, structName string) (specs []spec, err error) {
	var fieldCount int
	var s spec
	var subStructName string

	fieldCount = value.NumField()
	specs = make([]spec, fieldCount)

	for i := 0; i < fieldCount; i++ {
		s, err = buildSpecFromField(value.Field(i), value.Type().Field(i), structName)
		if err != nil {
			return nil, err
		}
		if s.Value.Kind() == reflect.Struct {
			s.Length = 0
			s.Repeat = 0
			subStructName = s.Value.Type().String()
			s.Children, err = buildSpecsFromStructValue(
				s.Value, subStructName)
			if err != nil {
				return nil, err
			}
		}
		specs[i] = s
	}
	return specs, nil
}

// Convert annotation on a structure into a specification for what
// should be read from a fixed field file.
func buildSpecs(structure interface{}) (specs []spec, err error) {
	var structValue, value reflect.Value
	var structType reflect.Type
	var structName string

	structValue = reflect.ValueOf(structure)
	structType = reflect.TypeOf(structure)
	structName = structType.String()

	value = structValue.Elem()
	specs, err = buildSpecsFromStructValue(value, structName)
	return specs, nil
}
