package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func BindInputs(op any, inputs map[any]any) error {
	if op == nil {
		return errors.New("op is nil")
	}

	opValue := reflect.ValueOf(op)
	if opValue.Kind() != reflect.Ptr || opValue.IsNil() {
		return errors.New("op is not a valid struct pointer")
	}

	opValue = opValue.Elem()
	opType := opValue.Type()

	argsValue, argsType := getTargetField(opValue, opType, "args")
	for i := 0; i < argsType.NumField(); i++ {
		field := argsType.Field(i)
		tags := strings.Split(field.Tag.Get("tf"), ",")
		for l := len(tags); l <= 3; l++ {
			tags = append(tags, "")
		}
		if tags[0] != "input" {
			continue
		}

		fieldName := tags[1]
		if fieldName == "" {
			fieldName = field.Name
		}
		value := inputs[fieldName]
		if value == nil {
			if tags[2] == "required" {
				return fmt.Errorf("missing required input: %s", fieldName)
			}
			continue
		}

		fieldValue := argsValue.Field(i)
		fieldType := fieldValue.Type()
		switch fieldType.Kind() {
		case reflect.Ptr:
			ptrValue := reflect.New(fieldType.Elem())
			ptrValue.Elem().Set(reflect.ValueOf(value))
			fieldValue.Set(ptrValue)
		case reflect.Slice:
			sliceValue := reflect.MakeSlice(fieldType, 0, 0)
			slice := reflect.ValueOf(value)
			for j := 0; j < slice.Len(); j++ {
				itemValue := reflect.New(fieldType.Elem()).Elem()
				itemValue.Set(slice.Index(j))
				sliceValue = reflect.Append(sliceValue, itemValue)
			}
			fieldValue.Set(sliceValue)
		default:
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func BindOutputs(op any) (map[any]any, error) {
	if op == nil {
		return nil, errors.New("op is nil")
	}

	opValue := reflect.ValueOf(op)
	if opValue.Kind() != reflect.Ptr || opValue.IsNil() {
		return nil, errors.New("op is not a valid struct pointer")
	}
	opValue = opValue.Elem()
	opType := opValue.Type()

	argsValue, argsType := getTargetField(opValue, opType, "args")
	outputs := make(map[any]any)
	for i := 0; i < argsType.NumField(); i++ {
		field := argsType.Field(i)
		tags := strings.Split(field.Tag.Get("tf"), ",")
		for l := len(tags); l <= 3; l++ {
			tags = append(tags, "")
		}
		if tags[0] != "output" {
			continue
		}

		fieldName := tags[1]
		if fieldName == "" {
			fieldName = field.Name
		}

		fieldValue := argsValue.Field(i)
		if !fieldValue.IsValid() {
			return nil, fmt.Errorf("invalid field: %s", fieldName)
		}
		if fieldValue.IsNil() {
			if tags[2] == "required" {
				return nil, fmt.Errorf("missing required output: %s", fieldName)
			}
			continue
		}
		if !fieldValue.CanInterface() {
			return nil, fmt.Errorf("can not interface field: %s", fieldName)
		}
		outputs[fieldName] = fieldValue.Interface()
	}
	return outputs, nil
}

func getTargetField(opValue reflect.Value, opType reflect.Type, tag string) (reflect.Value, reflect.Type) {
	var argsValue reflect.Value
	var argsType reflect.Type
	// Find the field with the specified tag
	for i := 0; i < opType.NumField(); i++ {
		field := opType.Field(i)
		if field.Tag.Get("tf") == tag {
			argsValue = opValue.Field(i)
			argsType = argsValue.Type()
			break
		}
	}
	if !argsValue.IsValid() {
		argsValue = opValue
		argsType = opType
	}
	return argsValue, argsType
}
