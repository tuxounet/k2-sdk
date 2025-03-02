package layouts

import "github.com/tuxounet/k2-sdk/kernel/storage/volumes/bases"

type TimeStream struct {
	bases.BaseLayout
}

func NewTimeStramLayout() *TimeStream {
	base := bases.NewBaseLayout()
	return &TimeStream{base}
}
