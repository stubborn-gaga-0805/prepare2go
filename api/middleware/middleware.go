package middleware

import (
	"context"
)

type Handler struct{}

func NewMiddleware(ctx context.Context) *Handler {
	return &Handler{}
}
