// Package scf implements serialization and deserialization of Simple Command Format,
// which is a binary format defined to be compatible with multiple layers of the ZigBee stack.
//
// SCF supports the following data types: 1-, 2-, 4- and 8-byte long signed and
// unsigned integers, as well as slices of such integers.
//
// Package SCF supports direct conversion between byte arrays and Go structs.
// The binary representation of a struct is defined as all its members in
// declaration order in little endian and with no padding. Slices are prefixed
// with a one byte count indicating the number of elements.
package scf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

// Validate checks that command has a type that is compatible with the
// requirements of this package and can be safely passed to Serialize and Parse.
// Command must be a struct and all fields must be supported data types (see package description).
// It returns an error if this is not the case.
func Validate(command interface{}) error {
	return ValidateType(reflect.TypeOf(command))
}

// Validate checks that commandType is compatible with the requirements of this
// package and can be safely passed to Serialize and Parse.
// Command must be a struct and all fields must be supported data types (see package description).
// It returns an error if this is not the case.
func ValidateType(commandType reflect.Type) error {
	if commandType.Kind() != reflect.Struct {
		return errors.New("command must be a struct")
	}

	for f := 0; f < commandType.NumField(); f++ {
		field := commandType.Field(f)
		if !isFieldTypeValid(field.Type) {
			return fmt.Errorf("command field %s.%s has invalid type", commandType.Name(), field.Name)
		}
	}

	return nil
}

func isFieldTypeValid(fieldType reflect.Type) bool {
	fieldKind := fieldType.Kind()
	if fieldKind == reflect.Slice {
		fieldKind = fieldType.Elem().Kind()
	}

	validKinds := []reflect.Kind{
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	}
	for _, kind := range validKinds {
		if fieldKind == kind {
			return true
		}
	}
	return false
}

// Serialize returns the binary representation of a command.
// Serialize panics if the type of command is not supported (use Validate to check).
// If a slice has more than 255 elements, it is truncated to the first 255 elements.
func Serialize(command interface{}) []byte {
	return SerializeValue(reflect.ValueOf(command))
}

// SerializeValue is the same as Serialize, except it accepts command already
// wrapped in a reflect.Value.
func SerializeValue(commandValue reflect.Value) (data []byte) {
	for f := 0; f < commandValue.NumField(); f++ {
		field := commandValue.Field(f)
		fieldKind := field.Kind()

		if fieldKind == reflect.Slice {
			length := field.Len()
			if length > math.MaxUint8 {
				length = math.MaxUint8
			}

			switch field.Type().Elem().Kind() {
			case reflect.Uint8:
				data = append(data, byte(length))
				data = append(data, field.Bytes()[:length]...)

			case reflect.Uint16:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]uint16)[:length] {
					start := len(data)
					data = append(data, 0, 0)
					binary.LittleEndian.PutUint16(data[start:], uint16(value))
				}

			case reflect.Uint32:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]uint32)[:length] {
					start := len(data)
					data = append(data, 0, 0, 0, 0)
					binary.LittleEndian.PutUint32(data[start:], uint32(value))
				}

			case reflect.Uint64:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]uint64)[:length] {
					start := len(data)
					data = append(data, 0, 0, 0, 0, 0, 0, 0, 0)
					binary.LittleEndian.PutUint64(data[start:], uint64(value))
				}

			case reflect.Int8:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]int8)[:length] {
					data = append(data, byte(value))
				}

			case reflect.Int16:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]int16)[:length] {
					start := len(data)
					data = append(data, 0, 0)
					binary.LittleEndian.PutUint16(data[start:], uint16(value))
				}

			case reflect.Int32:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]int32)[:length] {
					start := len(data)
					data = append(data, 0, 0, 0, 0)
					binary.LittleEndian.PutUint32(data[start:], uint32(value))
				}

			case reflect.Int64:
				data = append(data, byte(length))
				for _, value := range field.Interface().([]int64)[:length] {
					start := len(data)
					data = append(data, 0, 0, 0, 0, 0, 0, 0, 0)
					binary.LittleEndian.PutUint64(data[start:], uint64(value))
				}

			default:
				commandType := commandValue.Type()
				panic(fmt.Sprintf("serialization for command field not implemented: field=%s.%s type=%v kind=%v", commandType.Name(), commandType.Field(f).Name, field.Type(), fieldKind))
			}

			continue
		}

		var value uint64
		switch fieldKind {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = uint64(field.Int())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = field.Uint()
		default:
			commandType := commandValue.Type()
			panic(fmt.Sprintf("serialization for command field not implemented: field=%s.%s type=%v kind=%v", commandType.Name(), commandType.Field(f).Name, field.Type(), fieldKind))
		}

		switch fieldKind {
		case reflect.Int8, reflect.Uint8:
			data = append(data, byte(value))

		case reflect.Int16, reflect.Uint16:
			start := len(data)
			data = append(data, 0, 0)
			binary.LittleEndian.PutUint16(data[start:], uint16(value))

		case reflect.Int32, reflect.Uint32:
			start := len(data)
			data = append(data, 0, 0, 0, 0)
			binary.LittleEndian.PutUint32(data[start:], uint32(value))

		case reflect.Int64, reflect.Uint64:
			start := len(data)
			data = append(data, 0, 0, 0, 0, 0, 0, 0, 0)
			binary.LittleEndian.PutUint64(data[start:], uint64(value))
		}
	}

	return
}

var ErrInvalidData = errors.New("invalid data")

// Parse parses a command from the binary representation.
// command must be a pointer to a struct, as the results are written into it.
// Parse panics if the type of the struct is not supported (use Validate to check).
// If the command could not be parsed, ErrInvalidData is returned.
// Otherwise Parse returns the remaining part of data, that was not read into command.
func Parse(command interface{}, data []byte) ([]byte, error) {
	return ParseValue(reflect.ValueOf(command).Elem(), data)
}

// ParseValue is like Parse except it accepts a reflect.Value.
// In contrast to Parse, commandValue must wrap the struct directly (instead of being a pointer).
// Also, commandValue must be addressable, as the parsed values are written into it.
func ParseValue(commandValue reflect.Value, data []byte) ([]byte, error) {
	for f := 0; f < commandValue.NumField(); f++ {
		field := commandValue.Field(f)
		fieldKind := field.Kind()

		if fieldKind == reflect.Slice {
			if len(data) == 0 {
				return nil, ErrInvalidData
			}

			length := int(data[0])
			data = data[1:]

			elemKind := field.Type().Elem().Kind()

			switch elemKind {
			case reflect.Uint8, reflect.Int8:
				if len(data) < length {
					return nil, ErrInvalidData
				}

			case reflect.Uint16, reflect.Int16:
				if 2*len(data) < length {
					return nil, ErrInvalidData
				}

			case reflect.Uint32, reflect.Int32:
				if 4*len(data) < length {
					return nil, ErrInvalidData
				}

			case reflect.Uint64, reflect.Int64:
				if 8*len(data) < length {
					return nil, ErrInvalidData
				}

			default:
				commandType := commandValue.Type()
				panic(fmt.Sprintf("deserialization for command field not implemented: field=%s.%s type=%v kind=%v", commandType.Name(), commandType.Field(f).Name, field.Type(), fieldKind))
			}

			switch elemKind {
			case reflect.Uint8:
				fieldData := make([]uint8, length)
				copy(fieldData, data)
				data = data[length:]
				field.SetBytes(fieldData)

			case reflect.Uint16:
				fieldData := make([]uint16, length)
				for i := range fieldData {
					fieldData[i] = binary.LittleEndian.Uint16(data[2*i:])
				}
				data = data[2*length:]
				field.Set(reflect.ValueOf(fieldData))

			case reflect.Uint32:
				fieldData := make([]uint32, length)
				for i := range fieldData {
					fieldData[i] = binary.LittleEndian.Uint32(data[4*i:])
				}
				data = data[4*length:]
				field.Set(reflect.ValueOf(fieldData))

			case reflect.Uint64:
				fieldData := make([]uint64, length)
				for i := range fieldData {
					fieldData[i] = binary.LittleEndian.Uint64(data[8*i:])
				}
				data = data[8*length:]
				field.Set(reflect.ValueOf(fieldData))

			case reflect.Int8:
				fieldData := make([]int8, length)
				for i := range fieldData {
					fieldData[i] = int8(data[i])
				}
				data = data[length:]
				field.Set(reflect.ValueOf(fieldData))

			case reflect.Int16:
				fieldData := make([]int16, length)
				for i := range fieldData {
					fieldData[i] = int16(binary.LittleEndian.Uint16(data[2*i:]))
				}
				data = data[2*length:]
				field.Set(reflect.ValueOf(fieldData))

			case reflect.Int32:
				fieldData := make([]int32, length)
				for i := range fieldData {
					fieldData[i] = int32(binary.LittleEndian.Uint32(data[4*i:]))
				}
				data = data[4*length:]
				field.Set(reflect.ValueOf(fieldData))

			case reflect.Int64:
				fieldData := make([]int64, length)
				for i := range fieldData {
					fieldData[i] = int64(binary.LittleEndian.Uint64(data[8*i:]))
				}
				data = data[8*length:]
				field.Set(reflect.ValueOf(fieldData))
			}

			continue
		}

		length := int(field.Type().Size())
		if len(data) < length {
			return nil, ErrInvalidData
		}

		var value uint64
		switch fieldKind {
		case reflect.Int8, reflect.Uint8:
			value = uint64(data[0])

		case reflect.Int16, reflect.Uint16:
			value = uint64(binary.LittleEndian.Uint16(data))

		case reflect.Int32, reflect.Uint32:
			value = uint64(binary.LittleEndian.Uint32(data))

		case reflect.Int64, reflect.Uint64:
			value = uint64(binary.LittleEndian.Uint64(data))

		default:
			commandType := commandValue.Type()
			panic(fmt.Sprintf("deserialization for command field not implemented: field=%s.%s type=%v kind=%v", commandType.Name(), commandType.Field(f).Name, field.Type(), fieldKind))
		}

		data = data[length:]

		switch fieldKind {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(value))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(value)
		}
	}

	return data, nil
}
