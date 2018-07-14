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

func (gen IDGenerators) GenerateInt64ByType(v interface{}) int64 {
	return gen[reflect.TypeOf(v).Name()].GenerateInt64()
}

func New(maker func(typeString string) IDGenerator, values ...interface{}) IDGenerators {
	gens := make(map[string]IDGenerator)
	for _, v := range values {
		s := reflect.TypeOf(v).String()
		gens[s] = maker(s)
	}
	return gens
}
