package fixedfield

import (
	"fmt"
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

func (spec *readSpec) GetFieldName() string {
	return spec.FieldType.Name
}

func buildReadSpecs(structure interface{}) (readSpecs []readSpec, err error){
	var values reflect.Value 
	var spec readSpec
	var tag reflect.StructTag
	var length, repeat string

	values = reflect.ValueOf(structure).Elem()
	readSpecs = make([]readSpec, values.NumField())

	for i := 0; i < values.NumField(); i++ {
		spec = readSpecs[i]
		spec.FieldValue = values.Field(i)
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
		fmt.Printf(spec.String())
	}
	return readSpecs, nil
}

