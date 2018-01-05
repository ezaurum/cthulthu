package generators

type IDGenerator interface {
	Generate() string
	GenerateInt64() int64
}
