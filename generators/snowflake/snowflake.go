package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/ezaurum/cthulthu/generators"
)

var _ generators.IDGenerator = &snowflakeGenerator{}

func New(nodeNumber int64) generators.IDGenerator {
	node, err := snowflake.NewNode(nodeNumber)
	if nil != err {
		panic(err)
	}
	return &snowflakeGenerator{
		node: node,
	}
}

type snowflakeGenerator struct {
	node *snowflake.Node
}

func (g *snowflakeGenerator) Generate() string {
	s := g.node.Generate().String()
	return s
}

func (g *snowflakeGenerator) GenerateInt64() int64 {
	s := g.node.Generate().Int64()
	return s
}

func GetGenerators(nodeNumber int64, targets ...interface{}) generators.IDGenerators {
	gens := generators.New(func(_ string) generators.IDGenerator {
		return New(nodeNumber)
	}, targets...)
	return gens
}
