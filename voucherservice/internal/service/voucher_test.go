package service_test

import (
	"github.com/malkev1ch/observability/voucherservice/internal/service"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = service.Generate("B", 6)
	}
}
