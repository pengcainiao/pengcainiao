package collection

import (
	"github.com/jianfengye/collection"
)

type Kind uint

const (
	Invalid Kind = iota
	Int
	Int32
	Int64
	Uint
	Uint32
	Uint64
	Float32
	Float64
	String
	Object
)

func New(t Kind, data interface{}) collection.ICollection {
	switch t {
	case Int:
		return collection.NewIntCollection(data.([]int))
	case Int32:
		return collection.NewInt32Collection(data.([]int32))
	case Int64:
		return collection.NewInt64Collection(data.([]int64))
	case String:
		return collection.NewStrCollection(data.([]string))
	case Float64:
		return collection.NewFloat64Collection(data.([]float64))
	case Float32:
		return collection.NewFloat32Collection(data.([]float32))
	case Uint:
		return collection.NewUIntCollection(data.([]uint))
	case Uint32:
		return collection.NewUInt32Collection(data.([]uint32))
	case Uint64:
		return collection.NewUInt64Collection(data.([]uint64))
	case Object:
		return collection.NewObjCollection(data)
	default:
		return nil
	}
}
