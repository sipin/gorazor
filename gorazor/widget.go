package gorazor

type GorazorWidget interface {
	GetLabel() string
	GetValue() string
	GetName() string
	GetPlaceHolder() string
	GetType() string
}
