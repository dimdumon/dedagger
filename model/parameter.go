package model

import "reflect"

type Parameter struct {
	Type reflect.Type
}

func (p Parameter) String() string {
	return p.Type.Name()
}
