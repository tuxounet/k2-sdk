package paths

import (
	"path"
	"strings"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

const ServiceKey = "storage.paths"

type Service struct {
	bases.BaseKernelService
}

func NewService(k types.IKernel) types.IKernelService {

	base := bases.NewBaseKernelService(k, ServiceKey)
	return &Service{base}

}

func (c *Service) CominePath(parts ...string) string {
	return path.Join(parts...)
}

func (c *Service) SplitPath(parts ...string) []string {
	splitChar := "/"

	joined := path.Join(parts...)
	return strings.Split(joined, splitChar)

}

func (c *Service) GetBaseName(filePath string) string {
	return path.Base(filePath)
}

func (c *Service) GetDirName(filePath string) string {
	return path.Dir(filePath)
}

func (c *Service) GetExtName(filePath string) string {
	return path.Ext(filePath)
}
