package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Method struct {
	Name       string
	Parameters []*Parameter
	Value      reflect.Value
	StoreValue reflect.Value
}

func (m Method) String() string {
	parameterStrings := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		parameterStrings[i] = parameter.String()
	}
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(parameterStrings, ", "))
}
