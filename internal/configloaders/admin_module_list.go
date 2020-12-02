package configloaders

import "github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"

type AdminModuleList struct {
	IsSuper bool
	Modules []*systemconfigs.AdminModule
}
