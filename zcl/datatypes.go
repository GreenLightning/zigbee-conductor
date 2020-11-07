package zcl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type DataType byte

const (
	DataTypeNoData DataType = 0x00

	DataTypeData8  DataType = 0x08
	DataTypeData16 DataType = 0x09
	DataTypeData24 DataType = 0x0a
	DataTypeData32 DataType = 0x0b
	DataTypeData40 DataType = 0x0c
	DataTypeData48 DataType = 0x0d
	DataTypeData56 DataType = 0x0e
	DataTypeData64 DataType = 0x0f

	DataTypeBool DataType = 0x10

	DataTypeBitmap8  DataType = 0x18
	DataTypeBitmap16 DataType = 0x19
	DataTypeBitmap24 DataType = 0x1a
	DataTypeBitmap32 DataType = 0x1b
	DataTypeBitmap40 DataType = 0x1c
	DataTypeBitmap48 DataType = 0x1d
	DataTypeBitmap56 DataType = 0x1e
	DataTypeBitmap64 DataType = 0x1f

	DataTypeUint8  DataType = 0x20
	DataTypeUint16 DataType = 0x21
	DataTypeUint24 DataType = 0x22
	DataTypeUint32 DataType = 0x23
	DataTypeUint40 DataType = 0x24
	DataTypeUint48 DataType = 0x25
	DataTypeUint56 DataType = 0x26
	DataTypeUint64 DataType = 0x27

	DataTypeInt8  DataType = 0x28
	DataTypeInt16 DataType = 0x29
	DataTypeInt24 DataType = 0x2a
	DataTypeInt32 DataType = 0x2b
	DataTypeInt40 DataType = 0x2c
	DataTypeInt48 DataType = 0x2d
	DataTypeInt56 DataType = 0x2e
	DataTypeInt64 DataType = 0x2f

	DataTypeEnum8  DataType = 0x30
	DataTypeEnum16 DataType = 0x31

	DataTypeFloat16 DataType = 0x38
	DataTypeFloat32 DataType = 0x39
	DataTypeFloat64 DataType = 0x3a

	DataTypeOctetString         DataType = 0x41
	DataTypeCharacterString     DataType = 0x42
	DataTypeLongOctetString     DataType = 0x43
	DataTypeLongCharacterString DataType = 0x44

	DataTypeArray     DataType = 0x48
	DataTypeStructure DataType = 0x4c

	DataTypeSet DataType = 0x50
	DataTypeBag DataType = 0x51

	DataTypeTimeOfDay DataType = 0xe0
	DataTypeDate      DataType = 0xe1
	DataTypeUTCTime   DataType = 0xe2

	DataTypeClusterID   DataType = 0xe8
	DataTypeAttributeID DataType = 0xe9
	DataTypeBACnetOID   DataType = 0xea

	DataTypeIEEEAddress    DataType = 0xf0
	DataTypeSecurityKey128 DataType = 0xf1

	DataTypeUnknown DataType = 0xff
)

func (typ DataType) String() string {
	switch typ {
	case DataTypeNoData:
		return "nodata"

	case DataTypeData8:
		return "data8"
	case DataTypeData16:
		return "data16"
	case DataTypeData24:
		return "data24"
	case DataTypeData32:
		return "data32"
	case DataTypeData40:
		return "data40"
	case DataTypeData48:
		return "data48"
	case DataTypeData56:
		return "data56"
	case DataTypeData64:
		return "data64"

	case DataTypeBool:
		return "bool"

	case DataTypeBitmap8:
		return "bitmap8"
	case DataTypeBitmap16:
		return "bitmap16"
	case DataTypeBitmap24:
		return "bitmap24"
	case DataTypeBitmap32:
		return "bitmap32"
	case DataTypeBitmap40:
		return "bitmap40"
	case DataTypeBitmap48:
		return "bitmap48"
	case DataTypeBitmap56:
		return "bitmap56"
	case DataTypeBitmap64:
		return "bitmap64"

	case DataTypeUint8:
		return "uint8"
	case DataTypeUint16:
		return "uint16"
	case DataTypeUint24:
		return "uint24"
	case DataTypeUint32:
		return "uint32"
	case DataTypeUint40:
		return "uint40"
	case DataTypeUint48:
		return "uint48"
	case DataTypeUint56:
		return "uint56"
	case DataTypeUint64:
		return "uint64"

	case DataTypeInt8:
		return "int8"
	case DataTypeInt16:
		return "int16"
	case DataTypeInt24:
		return "int24"
	case DataTypeInt32:
		return "int32"
	case DataTypeInt40:
		return "int40"
	case DataTypeInt48:
		return "int48"
	case DataTypeInt56:
		return "int56"
	case DataTypeInt64:
		return "int64"

	case DataTypeEnum8:
		return "enum8"
	case DataTypeEnum16:
		return "enum16"

	case DataTypeFloat16:
		return "float16"
	case DataTypeFloat32:
		return "float32"
	case DataTypeFloat64:
		return "float64"

	case DataTypeOctetString:
		return "octstr"
	case DataTypeCharacterString:
		return "string"
	case DataTypeLongOctetString:
		return "octstr16"
	case DataTypeLongCharacterString:
		return "string16"

	case DataTypeArray:
		return "array"
	case DataTypeStructure:
		return "struct"

	case DataTypeSet:
		return "set"
	case DataTypeBag:
		return "bag"

	case DataTypeTimeOfDay:
		return "timeOfDay"
	case DataTypeDate:
		return "date"
	case DataTypeUTCTime:
		return "UTC"

	case DataTypeClusterID:
		return "clusterID"
	case DataTypeAttributeID:
		return "attributeID"
	case DataTypeBACnetOID:
		return "bacOID"

	case DataTypeIEEEAddress:
		return "EUI64"
	case DataTypeSecurityKey128:
		return "key128"

	case DataTypeUnknown:
		return "unknown"

	default:
		return fmt.Sprintf("Invalid(0x%0x)", byte(typ))
	}
}

var SizeInvalid = -1
var SizeDynamic = -2

func (typ DataType) SizeInBytes() int {
	switch typ {
	case DataTypeNoData:
		return 0

	case DataTypeData8, DataTypeBitmap8, DataTypeUint8, DataTypeInt8:
		return 1
	case DataTypeData16, DataTypeBitmap16, DataTypeUint16, DataTypeInt16:
		return 2
	case DataTypeData24, DataTypeBitmap24, DataTypeUint24, DataTypeInt24:
		return 3
	case DataTypeData32, DataTypeBitmap32, DataTypeUint32, DataTypeInt32:
		return 4
	case DataTypeData40, DataTypeBitmap40, DataTypeUint40, DataTypeInt40:
		return 5
	case DataTypeData48, DataTypeBitmap48, DataTypeUint48, DataTypeInt48:
		return 6
	case DataTypeData56, DataTypeBitmap56, DataTypeUint56, DataTypeInt56:
		return 7
	case DataTypeData64, DataTypeBitmap64, DataTypeUint64, DataTypeInt64:
		return 8

	case DataTypeBool:
		return 1

	case DataTypeEnum8:
		return 1
	case DataTypeEnum16:
		return 2

	case DataTypeFloat16:
		return 2
	case DataTypeFloat32:
		return 4
	case DataTypeFloat64:
		return 8

	case DataTypeOctetString, DataTypeCharacterString, DataTypeLongOctetString, DataTypeLongCharacterString:
		return SizeDynamic

	case DataTypeArray, DataTypeStructure, DataTypeSet, DataTypeBag:
		return SizeDynamic

	case DataTypeTimeOfDay:
		return 4
	case DataTypeDate:
		return 4
	case DataTypeUTCTime:
		return 4

	case DataTypeClusterID:
		return 2
	case DataTypeAttributeID:
		return 2
	case DataTypeBACnetOID:
		return 4

	case DataTypeIEEEAddress:
		return 8
	case DataTypeSecurityKey128:
		return 16

	case DataTypeUnknown:
		return 0

	default:
		return SizeInvalid
	}
}

var ErrNotEnoughData = errors.New("not enough data")
var ErrInvalidData = errors.New("invalid data")
var ErrNotImplemented = errors.New("not implemented")

func ParseValue(typ DataType, data []byte) (interface{}, []byte, error) {
	size := typ.SizeInBytes()
	if size >= 0 && len(data) < size {
		return nil, data, ErrNotEnoughData
	}

	switch typ {
	case DataTypeNoData:
		return nil, data, nil

	case DataTypeData8, DataTypeData16, DataTypeData24, DataTypeData32, DataTypeData40, DataTypeData48, DataTypeData56, DataTypeData64:
		value := data[0:size]
		return value, data[size:], nil

	case DataTypeBool:
		value := data[0]
		if value == 0 {
			return false, data[1:], nil
		} else if value == 1 {
			return true, data[1:], nil
		} else if value == 0xff {
			return nil, data[1:], nil
		} else {
			return nil, data[1:], ErrInvalidData
		}

	case DataTypeBitmap8:
		value := uint8(data[0])
		return value, data[1:], nil

	case DataTypeBitmap16:
		value := binary.LittleEndian.Uint16(data)
		return value, data[2:], nil

	case DataTypeBitmap24:
		var buffer [4]byte
		copy(buffer[:], data[:3])
		value := binary.LittleEndian.Uint32(buffer[:])
		return value, data[3:], nil

	case DataTypeBitmap32:
		value := binary.LittleEndian.Uint32(data)
		return value, data[4:], nil

	case DataTypeBitmap40, DataTypeBitmap48, DataTypeBitmap56:
		var buffer [8]byte
		copy(buffer[:], data[:size])
		value := binary.LittleEndian.Uint64(buffer[:])
		return value, data[size:], nil

	case DataTypeBitmap64:
		value := binary.LittleEndian.Uint64(data)
		return value, data[8:], nil

	case DataTypeUint8:
		value := uint8(data[0])
		if value == 0xff {
			return nil, data[1:], nil
		}
		return value, data[1:], nil

	case DataTypeUint16:
		value := binary.LittleEndian.Uint16(data)
		if value == 0xffff {
			return nil, data[2:], nil
		}
		return value, data[2:], nil

	case DataTypeUint24:
		var buffer [4]byte
		copy(buffer[:], data[:3])
		value := binary.LittleEndian.Uint32(buffer[:])
		if value == 0xffffff {
			return nil, data[3:], nil
		}
		return value, data[3:], nil

	case DataTypeUint32:
		value := binary.LittleEndian.Uint32(data)
		if value == 0xffff_ffff {
			return nil, data[4:], nil
		}
		return value, data[4:], nil

	case DataTypeUint40, DataTypeUint48, DataTypeUint56:
		invalid := true
		for i := 0; i < size; i++ {
			if data[i] != 0xff {
				invalid = false
				break
			}
		}
		if invalid {
			return nil, data[size:], nil
		}
		var buffer [8]byte
		copy(buffer[:], data[:size])
		value := binary.LittleEndian.Uint64(buffer[:])
		return value, data[size:], nil

	case DataTypeUint64:
		value := binary.LittleEndian.Uint64(data)
		if value == 0xffff_ffff_ffff_ffff {
			return nil, data[8:], nil
		}
		return value, data[8:], nil

	case DataTypeInt8:
		value := int8(data[0])
		if uint8(value) == 0x80 {
			return nil, data[1:], nil
		}
		return value, data[1:], nil

	case DataTypeInt16:
		value := int16(binary.LittleEndian.Uint16(data))
		if uint16(value) == 0x8000 {
			return nil, data[2:], nil
		}
		return value, data[2:], nil

	case DataTypeInt24:
		var buffer [4]byte
		copy(buffer[:], data[:3])
		extendSign(buffer[:], 3, 4)
		value := int32(binary.LittleEndian.Uint32(buffer[:]))
		if uint32(value) == 0xff800000 {
			return nil, data[3:], nil
		}
		return value, data[3:], nil

	case DataTypeInt32:
		value := int32(binary.LittleEndian.Uint32(data))
		if uint32(value) == 0x8000_0000 {
			return nil, data[4:], nil
		}
		return value, data[4:], nil

	case DataTypeInt40, DataTypeInt48, DataTypeInt56:
		invalid := true
		for i := 0; i < size-1; i++ {
			if data[i] != 0x00 {
				invalid = false
				break
			}
		}
		if invalid && data[size-1] == 0x80 {
			return nil, data[size:], nil
		}
		var buffer [8]byte
		copy(buffer[:], data[:size])
		extendSign(buffer[:], size, 8)
		value := int64(binary.LittleEndian.Uint64(buffer[:]))
		return value, data[size:], nil

	case DataTypeInt64:
		value := int64(binary.LittleEndian.Uint64(data))
		if uint64(value) == 0x8000_0000_0000_0000 {
			return nil, data[8:], nil
		}
		return value, data[8:], nil

	case DataTypeEnum8, DataTypeEnum16:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeFloat16:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeFloat32:
		bits := binary.LittleEndian.Uint32(data)
		value := math.Float32frombits(bits)
		if math.IsNaN(float64(value)) {
			return nil, data[4:], nil
		}
		return value, data[4:], nil

	case DataTypeFloat64:
		bits := binary.LittleEndian.Uint64(data)
		value := math.Float64frombits(bits)
		if math.IsNaN(value) {
			return nil, data[8:], nil
		}
		return value, data[8:], nil

	case DataTypeOctetString:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeCharacterString:
		if len(data) < 1 {
			return nil, data, ErrNotEnoughData
		}
		length := uint8(data[0])
		if length == 0xff {
			return nil, data[1:], nil
		}
		if len(data) < int(length)+1 {
			return nil, data, ErrNotEnoughData
		}
		value := string(data[1 : length+1])
		return value, data[length+1:], nil

	case DataTypeLongOctetString:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeLongCharacterString:
		if len(data) < 2 {
			return nil, data, ErrNotEnoughData
		}
		length := binary.LittleEndian.Uint16(data)
		if length == 0xffff {
			return nil, data[2:], nil
		}
		if len(data) < int(length)+2 {
			return nil, data, ErrNotEnoughData
		}
		value := string(data[2 : length+2])
		return value, data[length+2:], nil

	case DataTypeArray, DataTypeStructure:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeSet, DataTypeBag:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeTimeOfDay, DataTypeDate, DataTypeUTCTime:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeClusterID, DataTypeAttributeID, DataTypeBACnetOID:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeIEEEAddress, DataTypeSecurityKey128:
		return nil, data, fmt.Errorf("%w: %v", ErrNotImplemented, typ)

	case DataTypeUnknown:
		return nil, data, nil

	default:
		return nil, data, fmt.Errorf("invalid data type: 0x%0x", byte(typ))
	}
}

func extendSign(buffer []byte, validBytes, totalBytes int) {
	if buffer[validBytes-1]&0b1000_0000 != 0 {
		for i := validBytes; i < totalBytes; i++ {
			buffer[i] = 0xff
		}
	}
}
