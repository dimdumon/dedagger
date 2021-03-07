package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"

	"github.com/svarogg/dedagger/model"

	"github.com/gernest/utron/controller"
	"github.com/svarogg/dedagger/backend"
)

type Stores struct {
	controller.BaseController
	Routes  []string
	backend *backend.Backend
}

func (s *Stores) Home() {
	s.Ctx.Data["List"] = s.backend.Stores
	s.Ctx.Template = "stores"
	s.HTML(http.StatusOK)
}

func (s *Stores) Methods() {
	storeName := s.Ctx.Params["storeName"]
	var store *model.Store
	var ok bool
	if store, ok = s.backend.Stores[storeName]; !ok {
		s.HTML(http.StatusNotFound)
		return
	}

	s.Ctx.Data["Store"] = store
	s.Ctx.Template = "methods"
	s.HTML(http.StatusOK)
}

func (s *Stores) MethodInput() {
	storeName := s.Ctx.Params["storeName"]
	var store *model.Store
	var ok bool
	if store, ok = s.backend.Stores[storeName]; !ok {
		s.HTML(http.StatusNotFound)
		return
	}

	methodName := s.Ctx.Params["methodName"]
	var method *model.Method
	if method, ok = store.Methods[methodName]; !ok {
		s.HTML(http.StatusNotFound)
		return
	}

	s.Ctx.Data["Store"] = store
	s.Ctx.Data["Method"] = method
	s.Ctx.Template = "method"
	s.HTML(http.StatusOK)
}

func (s *Stores) MethodCall() {
	storeName := s.Ctx.Params["storeName"]
	var store *model.Store
	var ok bool
	if store, ok = s.backend.Stores[storeName]; !ok {
		s.HTML(http.StatusNotFound)
		return
	}

	methodName := s.Ctx.Params["methodName"]
	var method *model.Method
	if method, ok = store.Methods[methodName]; !ok {
		s.HTML(http.StatusNotFound)
		return
	}

	request := s.Ctx.Request()
	err := request.ParseForm()
	if err != nil {
		s.error(fmt.Errorf("error from parse form: %+v", err), http.StatusInternalServerError)
		return
	}

	parameterValues := make([]reflect.Value, len(method.Parameters))
	parameterStrings := make([]string, len(method.Parameters))
	for i, parameter := range method.Parameters {
		parameterString := request.FormValue(fmt.Sprintf("parameter%d", i))
		value, err := parseParameter(parameter, parameterString)
		if err != nil {
			s.error(fmt.Errorf("error from parseParameter: %+v", err), http.StatusBadRequest)
			return
		}
		parameterValues[i] = value
		parameterStrings[i] = parameterString
	}

	outputValues := s.backend.Call(method, parameterValues)
	outputs := make([]interface{}, len(outputValues))
	for i, outputValue := range outputValues {
		output, err := s.renderOutput(outputValue.Interface())
		if err != nil {
			s.error(fmt.Errorf("error from renderOutput: %+v", err), http.StatusBadRequest)
			return
		}

		outputs[i] = output
	}

	s.Ctx.Data["Store"] = store
	s.Ctx.Data["Method"] = method
	s.Ctx.Data["MethodWithParameters"] = fmt.Sprintf("%s(%s)", method, strings.Join(parameterStrings, ", "))
	s.Ctx.Data["Outputs"] = outputs
	s.Ctx.Template = "call_result"
	s.HTML(http.StatusOK)
}

func parseParameter(parameter *model.Parameter, valueString string) (reflect.Value, error) {
	switch parameter.Type.String() {
	case "*externalapi.DomainHash":
		hash, err := externalapi.NewDomainHashFromString(valueString)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(hash), nil
	default:
		valueInterface := reflect.New(parameter.Type).Interface()
		err := json.Unmarshal([]byte(valueString), valueInterface)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(valueInterface), nil
	}
}

func (s *Stores) error(err error, code int) {
	s.Ctx.Data["Error"] = fmt.Sprintf("%+v", err)
	s.Ctx.Template = "error"
	s.HTML(code)
}

func (s *Stores) renderOutput(output interface{}) (string, error) {
	bytes, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewStores(be *backend.Backend) func() controller.Controller {
	return func() controller.Controller {
		return &Stores{
			Routes: []string{
				"get;/;Home",
				"get;/{storeName}/methods;Methods",
				"get;/{storeName}/methods/{methodName};MethodInput",
				"post;/{storeName}/methods/{methodName};MethodCall",
			},
			backend: be,
		}
	}
}
