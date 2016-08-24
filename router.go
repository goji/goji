package goji

import (
	"github.com/weave-lab/goji/internal"
	"context"
)

type match struct {
	context.Context
	p Pattern
	h Handler
}

func (m match) Value(key interface{}) interface{} {
	switch key {
	case internal.Pattern:
		return m.p
	case internal.Handler:
		return m.h
	default:
		return m.Context.Value(key)
	}
}

var _ context.Context = match{}
