// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	ProcDir = "/proc"
)

func FindPidWithName(name string) int {
	// process name
	commFiles, err := filepath.Glob(ProcDir + "/*/comm")
	if err != nil {
		return 0
	}

	for _, commFile := range commFiles {
		data, err := os.ReadFile(commFile)
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(data)) == name {
			var pieces = strings.Split(commFile, "/")
			var pid = pieces[len(pieces)-2]
			pidInt, _ := strconv.Atoi(pid)
			return pidInt
		}
	}

	return 0
}
