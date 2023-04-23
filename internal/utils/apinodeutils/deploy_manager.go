package apinodeutils

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"regexp"
)

// DeployManager 节点部署文件管理器
// 如果节点部署文件有变化，需要重启API节点以便于生效
type DeployManager struct {
	dir string
}

// NewDeployManager 获取新节点部署文件管理器
func NewDeployManager() *DeployManager {
	var manager = &DeployManager{
		dir: Tea.Root + "/edge-api/deploy",
	}
	manager.LoadNodeFiles()
	manager.LoadNSNodeFiles()
	return manager
}

// LoadNodeFiles 加载所有边缘节点文件
func (this *DeployManager) LoadNodeFiles() []*DeployFile {
	var keyMap = map[string]*DeployFile{} // key => File

	var reg = regexp.MustCompile(`^edge-node-(\w+)-(\w+)-v([0-9.]+)\.zip$`)
	for _, file := range files.NewFile(this.dir).List() {
		var name = file.Name()
		if !reg.MatchString(name) {
			continue
		}
		var matches = reg.FindStringSubmatch(name)
		var osName = matches[1]
		var arch = matches[2]
		var version = matches[3]

		var key = osName + "_" + arch
		oldFile, ok := keyMap[key]
		if ok && stringutil.VersionCompare(oldFile.Version, version) > 0 {
			continue
		}
		keyMap[key] = &DeployFile{
			OS:      osName,
			Arch:    arch,
			Version: version,
			Path:    file.Path(),
		}
	}

	var result = []*DeployFile{}
	for _, v := range keyMap {
		result = append(result, v)
	}

	return result
}

// LoadNSNodeFiles 加载所有NS节点安装文件
func (this *DeployManager) LoadNSNodeFiles() []*DeployFile {
	var keyMap = map[string]*DeployFile{} // key => File

	var reg = regexp.MustCompile(`^edge-dns-(\w+)-(\w+)-v([0-9.]+)\.zip$`)
	for _, file := range files.NewFile(this.dir).List() {
		var name = file.Name()
		if !reg.MatchString(name) {
			continue
		}
		var matches = reg.FindStringSubmatch(name)
		var osName = matches[1]
		var arch = matches[2]
		var version = matches[3]

		var key = osName + "_" + arch
		oldFile, ok := keyMap[key]
		if ok && stringutil.VersionCompare(oldFile.Version, version) > 0 {
			continue
		}
		keyMap[key] = &DeployFile{
			OS:      osName,
			Arch:    arch,
			Version: version,
			Path:    file.Path(),
		}
	}

	var result = []*DeployFile{}
	for _, v := range keyMap {
		result = append(result, v)
	}

	return result
}
