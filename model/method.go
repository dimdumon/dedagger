package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Method struct {
	Name       string
	Value      reflect.Value
	Parameters []*Parameter
}

func (m Method) String() string {
	parameterStrings := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		parameterStrings[i] = parameter.String()
	}
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(parameterStrings, ", "))
}
