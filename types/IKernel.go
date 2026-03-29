package types

import (
	"context"
)

type IKernel interface {
	ILoggable
	IsUnsecure() bool
	IsForceCompute() bool
	GetApp() IApp
	GetRunDirectory() string
	GetRootContext() context.Context
	GetService(key KernelServiceContextKey) IKernelService
	SetService(service IKernelService)
	Init() error
	Register() error
	Start() error
	ListenAndServe() error
}
