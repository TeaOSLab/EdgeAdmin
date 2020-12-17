package rpc

import "context"

type ContextInterface interface {
	AdminContext() context.Context
}
