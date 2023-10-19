package handler

type Handler struct {
	*User
	*Voucher
}

func New(user *User, voucher *Voucher) *Handler {
	return &Handler{
		User:    user,
		Voucher: voucher,
	}
}
