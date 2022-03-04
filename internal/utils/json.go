// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package utils

import (
	"encoding/json"
	"reflect"
)

// JSONClone 使用JSON克隆对象
func JSONClone(v interface{}) (interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var nv = reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = json.Unmarshal(data, nv)
	if err != nil {
		return nil, err
	}

	return nv, nil
}
