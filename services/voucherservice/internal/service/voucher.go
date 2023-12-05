package service

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

	"github.com/malkev1ch/observability/services/voucherservice/internal/model"
)

type VoucherRepository interface {
	GetByUserID(ctx context.Context, id int64) (*model.Voucher, error)
	Create(ctx context.Context, user *model.Voucher) (*model.Voucher, error)
}

type Voucher struct {
	rps VoucherRepository
}

func NewVoucher(rps VoucherRepository) *Voucher {
	return &Voucher{rps: rps}
}

func (s *Voucher) GetByUserID(ctx context.Context, id int64) (*model.Voucher, error) {
	return s.rps.GetByUserID(ctx, id)
}

func (s *Voucher) Create(ctx context.Context, voucher *model.Voucher) (*model.Voucher, error) {
	const length = 6
	voucher.Value = Generate(strconv.Itoa(int(voucher.UserID)), length)
	return s.rps.Create(ctx, voucher)
}

// Generate generates random code according to prefix and length. Prefix should be ASCII symbols.
func Generate(prefix string, length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := &strings.Builder{}
	code.Grow(len(prefix) + length)
	code.WriteString(prefix)
	for i := 0; i < length; i++ {
		code.WriteRune(randomChar(charset))
	}

	return code.String()
}

// return random int in the range min...max.
func randomInt(min, max int) int {
	//nolint:gosec // it's ok to use math/rand here
	return min + rand.Intn(max-min)
}

// return random rune from charset.
func randomChar(cs string) rune {
	return rune(cs[randomInt(0, len(cs)-1)])
}
