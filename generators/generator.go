package generators

import (
	"reflect"
	"github.com/ezaurum/cthulthu/generators/snowflake"
)

type IDGenerator interface {
	Generate() string
	GenerateInt64() int64
}

type idGenerators struct {
	idGenerators map[string]IDGenerator
	nodeNumber   int64
}

func (gen *idGenerators) GenerateInt64(typeName string) int64 {
	return gen.idGenerators[typeName].GenerateInt64()
}

func (gen *idGenerators) Generate(typeName string) string {
	return gen.idGenerators[typeName].Generate()
}

func (gen *idGenerators) GenerateByType(v interface{}) string {
	return gen.idGenerators[reflect.TypeOf(v).Name()].Generate()
}

func (gen *idGenerators) GenerateInt64ByType(v interface{}) string {
	return gen.idGenerators[reflect.TypeOf(v).Name()].Generate()
}

func New(nodeNumber int64, values ...interface{}) *idGenerators {

	gens := make(map[string]IDGenerator)
	for _, v := range values {
		n := snowflake.New(nodeNumber)
		gens[reflect.TypeOf(v).Name()] = n
	}
	return &idGenerators{
		nodeNumber:   nodeNumber,
		idGenerators: gens,
	}
}
