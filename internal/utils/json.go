// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
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

// JSONDecodeConfig 解码并重新编码
// 是为了去除原有JSON中不需要的数据
func JSONDecodeConfig(data []byte, ptr any) (encodeJSON []byte, err error) {
	err = json.Unmarshal(data, ptr)
	if err != nil {
		return
	}

	encodeJSON, err = json.Marshal(ptr)
	if err != nil {
		return
	}

	// validate config
	if ptr != nil {
		config, ok := ptr.(interface {
			Init() error
		})
		if ok {
			initErr := config.Init()
			if initErr != nil {
				err = errors.New("validate config failed: " + initErr.Error())
			}
		}
	}

	return
}
