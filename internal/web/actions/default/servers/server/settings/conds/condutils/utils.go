package condutils

import (
	"encoding/json"
	"github.com/iwind/TeaGo/Tea"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
	"path/filepath"
)

type CondJSComponent struct {
	Type            string `json:"type"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Component       string `json:"component"`
	ParamsTitle     string `json:"paramsTitle"`
	IsRequest       bool   `json:"isRequest"`
	CaseInsensitive bool   `json:"caseInsensitive"`
}

// ReadAllAvailableCondTypes 读取所有可用的条件
func ReadAllAvailableCondTypes() []*CondJSComponent {
	result := []*CondJSComponent{}

	dir := Tea.Root + "/web/"
	if Tea.IsTesting() {
		dir = filepath.Dir(Tea.Root) + "/web"
	}
	dir += "/public/js/conds/"
	jsonFiles := files.NewFile(dir).List()
	for _, file := range jsonFiles {
		if file.Ext() == ".json" {
			data, err := file.ReadAll()
			if err != nil {
				logs.Println("[COND]read data from json file: " + err.Error())
				continue
			}

			c := []*CondJSComponent{}
			err = json.Unmarshal(data, &c)
			if err != nil {
				logs.Println("[COND]decode json failed: " + err.Error())
				continue
			}
			result = append(result, c...)
		}
	}

	return result
}
