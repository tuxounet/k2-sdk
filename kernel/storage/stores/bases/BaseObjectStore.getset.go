package bases

import "github.com/tuxounet/k2-sdk/types"

func (s *BaseObjectStore[R]) GetLogger() types.ILogger {
	return s.log
}

func (s *BaseObjectStore[R]) GetKernel() types.IKernel {
	return s.kernel
}
