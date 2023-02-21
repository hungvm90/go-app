package internal

import (
	"context"
	"log"
	"sync"
)

type AppContext struct {
	Wg *sync.WaitGroup
	context.Context
}

type RequestContext struct {
	Logger *log.Logger
	context.Context
}
