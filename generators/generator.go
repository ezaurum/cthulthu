package generators

import (
	"reflect"
)

type IDGenerator interface {
	Generate() string
	GenerateInt64() int64
}

type IDGenerators struct {
	generators       map[string]IDGenerator
	defaultGenerator IDGenerator
}

func (gen IDGenerators) GenerateInt64(typeName string) int64 {
	if generator, b := gen.generators[typeName]; b {
		return generator.GenerateInt64()
	} else {
		return gen.defaultGenerator.GenerateInt64()
	}
}

func (gen IDGenerators) Generate(typeName string) string {
	if generator, b := gen.generators[typeName]; b {
		return generator.Generate()
	} else {
		return gen.defaultGenerator.Generate()
	}
}

func (gen IDGenerators) GenerateByType(v interface{}) string {
	return gen.Generate(reflect.TypeOf(v).Name())
}

func (gen IDGenerators) GenerateInt64ByType(v interface{}) int64 {
	return gen.GenerateInt64(reflect.TypeOf(v).Name())
}

func (gen *IDGenerators) Add(v interface{}, maker GeneratorMakerFunc) {
	s := reflect.TypeOf(v).String()
	gen.generators[s] = maker("default")
}

func New(maker GeneratorMakerFunc, values ...interface{}) IDGenerators {
	gens := make(map[string]IDGenerator)
	for _, v := range values {
		s := reflect.TypeOf(v).String()
		gens[s] = maker(s)
	}
	defaultGen := maker("default")
	return IDGenerators{
		generators:       gens,
		defaultGenerator: defaultGen,
	}
}

type GeneratorMakerFunc func(typeString string) IDGenerator
