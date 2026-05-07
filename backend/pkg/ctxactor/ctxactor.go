package ctxactor

import "context"

type Actor struct {
	UserID int64
	NIK    string
	Name   string
}

type keyType struct{}

var key keyType

// Simpan actor ke context
func WithActor(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, key, a)
}

// Ambil actor dari context
func From(ctx context.Context) (Actor, bool) {
	a, ok := ctx.Value(key).(Actor)
	return a, ok
}
