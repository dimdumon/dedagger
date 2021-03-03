package model

import (
	"reflect"
)

type Store struct {
	Name    string
	Value   reflect.Value
	Typeof  reflect.Type
	Methods []*Method
}

func (s Store) String() string {
	return s.Name
}
