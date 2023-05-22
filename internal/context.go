package internal

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type AppContext struct {
	Wg *sync.WaitGroup
	context.Context
}

type RequestContext struct {
	Logger *zerolog.Logger
	context.Context
}

func (r *RequestContext) WithLogger(logger zerolog.Logger) {
	temp := logger
	r.Logger = &temp
}

func NewRequestContext(ctx context.Context) RequestContext {
	return RequestContext{Logger: &log.Logger, Context: ctx}
}
