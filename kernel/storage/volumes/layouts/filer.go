package layouts

import "github.com/tuxounet/k2-sdk/kernel/storage/volumes/bases"

type Filer struct {
	bases.BaseLayout
}

func NewFilerLayout() *Filer {
	base := bases.NewBaseLayout()
	return &Filer{base}
}
