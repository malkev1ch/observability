package service_test

import (
	"testing"

	"github.com/malkev1ch/observability/services/voucherservice/internal/service"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	voucher := service.Generate("test", 3)
	require.Len(t, voucher, 7)
	require.Contains(t, voucher, "test")
}
