// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	executils "github.com/TeaOSLab/EdgeAdmin/internal/utils/exec"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type UpgradeFileWriter struct {
	rawWriter io.Writer
	written   int64
}

func NewUpgradeFileWriter(rawWriter io.Writer) *UpgradeFileWriter {
	return &UpgradeFileWriter{rawWriter: rawWriter}
}

func (this *UpgradeFileWriter) Write(p []byte) (n int, err error) {
	n, err = this.rawWriter.Write(p)
	this.written += int64(n)
	return
}

func (this *UpgradeFileWriter) TotalWritten() int64 {
	return this.written
}

type UpgradeManager struct {
	client *http.Client

	component string

	newVersion    string
	contentLength int64
	isDownloading bool
	writer        *UpgradeFileWriter
	body          io.ReadCloser
	isCancelled   bool

	downloadURL string
}

func NewUpgradeManager(component string, downloadURL string) *UpgradeManager {
	return &UpgradeManager{
		component:   component,
		downloadURL: downloadURL,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       30 * time.Minute,
		},
	}
}

func (this *UpgradeManager) Start() error {
	if this.isDownloading {
		return errors.New("another process is running")
	}

	this.isDownloading = true

	defer func() {
		this.client.CloseIdleConnections()
		this.isDownloading = false
	}()

	// 检查unzip
	unzipExe, _ := executils.LookPath("unzip")

	// 检查cp
	cpExe, _ := executils.LookPath("cp")
	if len(cpExe) == 0 {
		return errors.New("can not find 'cp' command")
	}

	// 检查新版本
	var downloadURL = this.downloadURL
	if len(downloadURL) == 0 {
		var url = teaconst.UpdatesURL
		var osName = runtime.GOOS
		if Tea.IsTesting() && osName == "darwin" {
			osName = "linux"
		}
		url = strings.ReplaceAll(url, "${os}", osName)
		url = strings.ReplaceAll(url, "${arch}", runtime.GOARCH)
		url = strings.ReplaceAll(url, "${version}", teaconst.Version)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("create url request failed: %w", err)
		}
		req.Header.Set("User-Agent", "Edge-Admin/"+teaconst.Version)

		resp, err := this.client.Do(req)
		if err != nil {
			return fmt.Errorf("read latest version failed: %w", err)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			return errors.New("read latest version failed: invalid response code '" + types.String(resp.StatusCode) + "'")
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read latest version failed: %w", err)
		}

		var m = maps.Map{}
		err = json.Unmarshal(data, &m)
		if err != nil {
			return fmt.Errorf("invalid response data: %w, origin data: %s", err, string(data))
		}

		var code = m.GetInt("code")
		if code != 200 {
			return errors.New(m.GetString("message"))
		}

		var dataMap = m.GetMap("data")
		var downloadHost = dataMap.GetString("host")
		var versions = dataMap.GetSlice("versions")
		var downloadPath = ""
		for _, component := range versions {
			var componentMap = maps.NewMap(component)
			if componentMap.Has("version") {
				if componentMap.GetString("code") == this.component {
					var version = componentMap.GetString("version")
					if stringutil.VersionCompare(version, teaconst.Version) > 0 {
						this.newVersion = version
						downloadPath = componentMap.GetString("url")
						break
					}
				}
			}
		}

		if len(downloadPath) == 0 {
			return errors.New("no latest version to download")
		}

		downloadURL = downloadHost + downloadPath
	}

	{
		req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
		if err != nil {
			return fmt.Errorf("create download request failed: %w", err)
		}
		req.Header.Set("User-Agent", "Edge-Admin/"+teaconst.Version)

		resp, err := this.client.Do(req)
		if err != nil {
			return fmt.Errorf("download failed: '%s': %w", downloadURL, err)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			return errors.New("download failed: " + downloadURL + ": invalid response code '" + types.String(resp.StatusCode) + "'")
		}

		this.contentLength = resp.ContentLength
		this.body = resp.Body

		// download to tmp
		var tmpDir = os.TempDir()
		var filename = filepath.Base(downloadURL)

		var destFile = tmpDir + "/" + filename
		_ = os.Remove(destFile)

		fp, err := os.Create(destFile)
		if err != nil {
			return fmt.Errorf("create file failed: %w", err)
		}

		defer func() {
			// 删除安装文件
			_ = os.Remove(destFile)
		}()

		this.writer = NewUpgradeFileWriter(fp)

		_, err = io.Copy(this.writer, resp.Body)
		if err != nil {
			_ = fp.Close()
			if this.isCancelled {
				return nil
			}
			return fmt.Errorf("download failed: %w", err)
		}

		_ = fp.Close()

		// unzip
		var unzipDir = tmpDir + "/edge-" + this.component + "-tmp"
		stat, err := os.Stat(unzipDir)
		if err == nil && stat.IsDir() {
			err = os.RemoveAll(unzipDir)
			if err != nil {
				return fmt.Errorf("remove old dir '%s' failed: %w", unzipDir, err)
			}
		}

		if len(unzipExe) > 0 {
			var unzipCmd = exec.Command(unzipExe, "-q", "-o", destFile, "-d", unzipDir)
			var unzipStderr = &bytes.Buffer{}
			unzipCmd.Stderr = unzipStderr
			err = unzipCmd.Run()
			if err != nil {
				return fmt.Errorf("unzip installation file failed: %w: %s", err, unzipStderr.String())
			}
		} else {
			var unzipCmd = &Unzip{
				zipFile:   destFile,
				targetDir: unzipDir,
			}
			err = unzipCmd.Run()
			if err != nil {
				return fmt.Errorf("unzip installation file failed: %w", err)
			}
		}

		installationFiles, err := filepath.Glob(unzipDir + "/edge-" + this.component + "/*")
		if err != nil {
			return fmt.Errorf("lookup installation files failed: %w", err)
		}

		// cp to target dir
		currentExe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("reveal current executable file path failed: %w", err)
		}
		var targetDir = filepath.Dir(filepath.Dir(currentExe))
		if !Tea.IsTesting() {
			for _, installationFile := range installationFiles {
				var cpCmd = exec.Command(cpExe, "-R", "-f", installationFile, targetDir)
				var cpStderr = &bytes.Buffer{}
				cpCmd.Stderr = cpStderr
				err = cpCmd.Run()
				if err != nil {
					return errors.New("overwrite installation files failed: '" + cpCmd.String() + "': " + cpStderr.String())
				}
			}
		}

		// remove tmp
		_ = os.RemoveAll(unzipDir)
	}

	return nil
}

func (this *UpgradeManager) IsDownloading() bool {
	return this.isDownloading
}

func (this *UpgradeManager) Progress() float32 {
	if this.contentLength <= 0 {
		return -1
	}
	if this.writer == nil {
		return -1
	}
	return float32(this.writer.TotalWritten()) / float32(this.contentLength)
}

func (this *UpgradeManager) NewVersion() string {
	return this.newVersion
}

func (this *UpgradeManager) Cancel() error {
	this.isCancelled = true
	this.isDownloading = false

	if this.body != nil {
		_ = this.body.Close()
	}
	return nil
}
