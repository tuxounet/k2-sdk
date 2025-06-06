package types

type IBaseObjectStore[R any] interface {
	HasValue() (bool, error)
	GetValue() (*R, error)
	SetValue(value R) error
	Nuke() error
}
