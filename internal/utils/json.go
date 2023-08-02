// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package utils

import (
	"bytes"
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

// JSONIsNull 判断JSON数据是否为null
func JSONIsNull(jsonData []byte) bool {
	return len(jsonData) == 0 || bytes.Equal(jsonData, []byte("null"))
}
