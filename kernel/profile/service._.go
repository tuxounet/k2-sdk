package profile

import (
	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "profile"

type ProfileService struct {
	bases.BaseKernelService
}

func NewService(k types.IKernel) types.IKernelService {

	base := bases.NewBaseKernelService(k, ServiceKey)
	return &ProfileService{base}

}
