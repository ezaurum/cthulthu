package generators

import (
	"reflect"
)

type IDGenerator interface {
	Generate() string
	GenerateInt64() int64
}

type IDGenerators map[string]IDGenerator

func (gen IDGenerators) GenerateInt64(typeName string) int64 {
	return gen[typeName].GenerateInt64()
}

func (gen IDGenerators) Generate(typeName string) string {
	return gen[typeName].Generate()
}

func (gen IDGenerators) GenerateByType(v interface{}) string {
	return gen[reflect.TypeOf(v).Name()].Generate()
}

func (gen IDGenerators) GenerateInt64ByType(v interface{}) string {
	return gen[reflect.TypeOf(v).Name()].Generate()
}

func New(maker func() IDGenerator, values ...interface{}) IDGenerators {
	gens := make(map[string]IDGenerator)
	for _, v := range values {
		n := maker()
		gens[reflect.TypeOf(v).Name()] = n
	}
	return gens
}
