package encrypt

import (
	"errors"
	"reflect"
)

var methods = map[string]reflect.Type{
	"raw":         reflect.TypeOf(new(RawMethod)).Elem(),
	"aes-128-cfb": reflect.TypeOf(new(AES128CFBMethod)).Elem(),
	"aes-192-cfb": reflect.TypeOf(new(AES192CFBMethod)).Elem(),
	"aes-256-cfb": reflect.TypeOf(new(AES256CFBMethod)).Elem(),
}

func NewMethodInstance(method string, key string, iv string) (MethodInterface, error) {
	valueType, ok := methods[method]
	if !ok {
		return nil, errors.New("method '" + method + "' not found")
	}
	instance, ok := reflect.New(valueType).Interface().(MethodInterface)
	if !ok {
		return nil, errors.New("method '" + method + "' must implement MethodInterface")
	}
	err := instance.Init([]byte(key), []byte(iv))
	return instance, err
}

func RecoverMethodPanic(err interface{}) error {
	if err != nil {
		s, ok := err.(string)
		if ok {
			return errors.New(s)
		}

		e, ok := err.(error)
		if ok {
			return e
		}

		return errors.New("unknown error")
	}
	return nil
}
