package service_test

import (
	"github.com/malkev1ch/observability/voucherservice/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	voucher := service.Generate("test", 3)
	require.Len(t, voucher, 7)
	require.Contains(t, voucher, "test")
}
