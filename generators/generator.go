package generators

import (
	"reflect"
)

type IDGenerator interface {
	Generate() string
	GenerateInt64() int64
}

type IDGenerators struct {
	IDGenerators map[string]IDGenerator
}

func (gen *IDGenerators) GenerateInt64(typeName string) int64 {
	return gen.IDGenerators[typeName].GenerateInt64()
}

func (gen *IDGenerators) Generate(typeName string) string {
	return gen.IDGenerators[typeName].Generate()
}

func (gen *IDGenerators) GenerateByType(v interface{}) string {
	return gen.IDGenerators[reflect.TypeOf(v).Name()].Generate()
}

func (gen *IDGenerators) GenerateInt64ByType(v interface{}) string {
	return gen.IDGenerators[reflect.TypeOf(v).Name()].Generate()
}

func New(generator IDGenerator, values ...interface{}) *IDGenerators {

	gens := make(map[string]IDGenerator)
	for _, v := range values {
		n := generator
		gens[reflect.TypeOf(v).Name()] = n
	}
	return &IDGenerators{
		IDGenerators: gens,
	}
}
