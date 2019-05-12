package barcode_snowflake

import (
	"testing"
)

func BenchmarkGenerate(b *testing.B) {

	node, _ := New(0)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}
