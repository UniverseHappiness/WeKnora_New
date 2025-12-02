package test

import (
	"context"
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
)

func TestContextKeyString(t *testing.T) {
	cases := []struct {
		key  types.ContextKey
		want string
	}{
		{types.TenantIDContextKey, "TenantID"},
		{types.TenantInfoContextKey, "TenantInfo"},
		{types.RequestIDContextKey, "RequestID"},
		{types.LoggerContextKey, "Logger"},
	}
	for _, c := range cases {
		if got := c.key.String(); got != c.want {
			t.Fatalf("String() = %q, want %q", got, c.want)
		}
	}
}

func TestContextKeyIsolation(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.TenantIDContextKey, uint(123))
	if v := ctx.Value(types.TenantIDContextKey); v != uint(123) {
		t.Fatalf("Value(TenantIDContextKey) = %v, want %v", v, uint(123))
	}
	ctx = context.WithValue(ctx, "TenantID", uint(456))
	if v := ctx.Value(types.TenantIDContextKey); v != uint(123) {
		t.Fatalf("typed key must not collide with string key; got %v", v)
	}
	if v := ctx.Value("TenantID"); v != uint(456) {
		t.Fatalf("string key retrieval failed, got %v", v)
	}
}
