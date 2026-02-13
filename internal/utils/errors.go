package utils

import "errors"

var (
	// Client errors
	ErrValidation      = errors.New("validation failed") // input tidak valid
	ErrUnauthorized    = errors.New("unauthorized")      // tidak ada login / token invalid
	ErrForbidden       = errors.New("forbidden")         // tidak punya akses
	ErrNotFound        = errors.New("data not found")    // resource tidak ditemukan
	ErrConflict        = errors.New("conflict")          // sudah ada (duplicate)
	ErrTooManyRequest  = errors.New("too many requests") // rate limit / throttle
	ErrInvalidPassword = errors.New("invalid password")  // rate limit / throttle
	ErrInvalidEmail    = errors.New("invalid email")     // rate limit / throttle

	// Server errors
	ErrInternal    = errors.New("internal server error") // kesalahan server
	ErrUnavailable = errors.New("service unavailable")   // service down / maintenance
	ErrTimeout     = errors.New("request timeout")       // koneksi lama / gagal

	// Payment / third-party errors
	ErrPayment          = errors.New("payment error")           // general payment error
	ErrPaymentDeclined  = errors.New("payment declined")        // ditolak oleh provider
	ErrPaymentExpired   = errors.New("payment expired")         // sudah lewat waktu bayar
	ErrPaymentCancelled = errors.New("payment cancelled")       // dibatalkan user / sistem
	ErrPaymentPending   = errors.New("payment pending")         // masih menunggu pembayaran
	ErrIntegration      = errors.New("integration error")       // error komunikasi dengan 3rd party
	ErrInvalidSignature = errors.New("invalid signature error") // signature tidak cocok (security)
)
