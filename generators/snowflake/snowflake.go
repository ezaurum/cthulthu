package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/ezaurum/cthulthu/generators"
)

var _ generators.IDGenerator = snowflakeGenerator{}

func New(nodeNumber int64) generators.IDGenerator {
	node, err := snowflake.NewNode(nodeNumber)
	if nil != err {
		panic(err)
	}
	return snowflakeGenerator{
		node: node,
	}
}

type snowflakeGenerator struct {
	node *snowflake.Node
}

func (g snowflakeGenerator) Generate() string {
	return g.node.Generate().String()
}
