package internal

import (
	"context"
	"github.com/rs/zerolog"
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
