// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package apinodeutils

import (
	"compress/gzip"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/types"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type Progress struct {
	Percent float64
}

type Upgrader struct {
	progress  *Progress
	apiExe    string
	apiNodeId int64
}

func NewUpgrader(apiNodeId int64) *Upgrader {
	return &Upgrader{
		apiExe:    apiExe(),
		progress:  &Progress{Percent: 0},
		apiNodeId: apiNodeId,
	}
}

func (this *Upgrader) Upgrade() error {
	sharedClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	apiNodeResp, err := sharedClient.APINodeRPC().FindEnabledAPINode(sharedClient.Context(0), &pb.FindEnabledAPINodeRequest{ApiNodeId: this.apiNodeId})
	if err != nil {
		return err
	}
	var apiNode = apiNodeResp.ApiNode
	if apiNode == nil {
		return errors.New("could not find api node with id '" + types.String(this.apiNodeId) + "'")
	}

	apiConfig, err := configs.LoadAPIConfig()
	if err != nil {
		return err
	}
	var newAPIConfig = apiConfig.Clone()
	newAPIConfig.RPCEndpoints = apiNode.AccessAddrs

	rpcClient, err := rpc.NewRPCClient(newAPIConfig, false)
	if err != nil {
		return err
	}

	// 升级边缘节点
	err = this.upgradeNodes(sharedClient.Context(0), rpcClient)
	if err != nil {
		return err
	}

	// 升级NS节点
	err = this.upgradeNSNodes(sharedClient.Context(0), rpcClient)
	if err != nil {
		return err
	}

	// 升级API节点
	err = this.upgradeAPINode(sharedClient.Context(0), rpcClient)
	if err != nil {
		return fmt.Errorf("upgrade api node failed: %w", err)
	}

	return nil
}

// Progress 查看升级进程
func (this *Upgrader) Progress() *Progress {
	return this.progress
}

// 升级API节点
func (this *Upgrader) upgradeAPINode(ctx context.Context, rpcClient *rpc.RPCClient) error {
	versionResp, err := rpcClient.APINodeRPC().FindCurrentAPINodeVersion(ctx, &pb.FindCurrentAPINodeVersionRequest{})
	if err != nil {
		return err
	}
	if !Tea.IsTesting() /** 开发环境下允许突破此限制方便测试 **/ &&
		(stringutil.VersionCompare(versionResp.Version, "0.6.4" /** 从0.6.4开始支持 **/) < 0 || versionResp.Os != runtime.GOOS || versionResp.Arch != runtime.GOARCH) {
		return errors.New("could not upgrade api node v" + versionResp.Version + "/" + versionResp.Os + "/" + versionResp.Arch)
	}

	// 检查本地文件版本
	canUpgrade, reason := CanUpgrade(versionResp.Version, versionResp.Os, versionResp.Arch)
	if !canUpgrade {
		return errors.New(reason)
	}

	localVersion, err := lookupLocalVersion()
	if err != nil {
		return fmt.Errorf("lookup version failed: %w", err)
	}

	// 检查要升级的文件
	var gzFile = this.apiExe + "." + localVersion + ".gz"

	gzReader, err := os.Open(gzFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = func() error {
			// 压缩文件
			exeReader, err := os.Open(this.apiExe)
			if err != nil {
				return err
			}
			defer func() {
				_ = exeReader.Close()
			}()
			var tmpGzFile = gzFile + ".tmp"
			gzFileWriter, err := os.OpenFile(tmpGzFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
			var gzWriter = gzip.NewWriter(gzFileWriter)
			defer func() {
				_ = gzWriter.Close()
				_ = gzFileWriter.Close()

				_ = os.Rename(tmpGzFile, gzFile)
			}()
			_, err = io.Copy(gzWriter, exeReader)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
		gzReader, err = os.Open(gzFile)
		if err != nil {
			return err
		}
	}

	defer func() {
		_ = gzReader.Close()
	}()

	// 开始上传
	var hash = md5.New()
	var buf = make([]byte, 128*4096)
	var isFirst = true
	stat, err := gzReader.Stat()
	if err != nil {
		return err
	}
	var totalSize = stat.Size()
	if totalSize == 0 {
		_ = gzReader.Close()
		_ = os.Remove(gzFile)
		return errors.New("invalid gz file")
	}

	var uploadedSize int64 = 0
	for {
		n, err := gzReader.Read(buf)
		if n > 0 {
			// 计算Hash
			hash.Write(buf[:n])

			// 上传
			_, uploadErr := rpcClient.APINodeRPC().UploadAPINodeFile(rpcClient.Context(0), &pb.UploadAPINodeFileRequest{
				Filename:     filepath.Base(this.apiExe),
				Sum:          "",
				ChunkData:    buf[:n],
				IsFirstChunk: isFirst,
				IsLastChunk:  false,
			})
			if uploadErr != nil {
				return uploadErr
			}

			// 进度
			uploadedSize += int64(n)
			this.progress = &Progress{Percent: float64(uploadedSize*100) / float64(totalSize)}
		}
		if isFirst {
			isFirst = false
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			if err == io.EOF {
				_, uploadErr := rpcClient.APINodeRPC().UploadAPINodeFile(rpcClient.Context(0), &pb.UploadAPINodeFileRequest{
					Filename:     filepath.Base(this.apiExe),
					Sum:          fmt.Sprintf("%x", hash.Sum(nil)),
					ChunkData:    buf[:n],
					IsFirstChunk: isFirst,
					IsLastChunk:  true,
				})
				if uploadErr != nil {
					return uploadErr
				}

				break
			}
		}
	}

	return nil
}

// 升级边缘节点
func (this *Upgrader) upgradeNodes(ctx context.Context, rpcClient *rpc.RPCClient) error {
	// 本地的
	var manager = NewDeployManager()
	var localFileMap = map[string]*DeployFile{} // os_arch => *DeployFile
	for _, deployFile := range manager.LoadNodeFiles() {
		localFileMap[deployFile.OS+"_"+deployFile.Arch] = deployFile
	}

	remoteFilesResp, err := rpcClient.APINodeRPC().FindLatestDeployFiles(ctx, &pb.FindLatestDeployFilesRequest{})
	if err != nil {
		return err
	}

	var remoteFileMap = map[string]*pb.FindLatestDeployFilesResponse_DeployFile{} // os_arch => *DeployFile
	for _, nodeFile := range remoteFilesResp.NodeDeployFiles {
		remoteFileMap[nodeFile.Os+"_"+nodeFile.Arch] = nodeFile
	}

	// 对比版本
	for key, deployFile := range localFileMap {
		remoteDeployFile, ok := remoteFileMap[key]
		if !ok || stringutil.VersionCompare(remoteDeployFile.Version, deployFile.Version) < 0 {
			err = this.uploadNodeDeployFile(ctx, rpcClient, deployFile.Path)
			if err != nil {
				return fmt.Errorf("upload deploy file '%s' failed: %w", filepath.Base(deployFile.Path), err)
			}
		}
	}

	return nil
}

// 升级NS节点
func (this *Upgrader) upgradeNSNodes(ctx context.Context, rpcClient *rpc.RPCClient) error {
	// 本地的
	var manager = NewDeployManager()
	var localFileMap = map[string]*DeployFile{} // os_arch => *DeployFile
	for _, deployFile := range manager.LoadNSNodeFiles() {
		localFileMap[deployFile.OS+"_"+deployFile.Arch] = deployFile
	}

	remoteFilesResp, err := rpcClient.APINodeRPC().FindLatestDeployFiles(ctx, &pb.FindLatestDeployFilesRequest{})
	if err != nil {
		return err
	}

	var remoteFileMap = map[string]*pb.FindLatestDeployFilesResponse_DeployFile{} // os_arch => *DeployFile
	for _, nodeFile := range remoteFilesResp.NsNodeDeployFiles {
		remoteFileMap[nodeFile.Os+"_"+nodeFile.Arch] = nodeFile
	}

	// 对比版本
	for key, deployFile := range localFileMap {
		remoteDeployFile, ok := remoteFileMap[key]
		if !ok || stringutil.VersionCompare(remoteDeployFile.Version, deployFile.Version) < 0 {
			err = this.uploadNodeDeployFile(ctx, rpcClient, deployFile.Path)
			if err != nil {
				return fmt.Errorf("upload deploy file '%s' failed: %w", filepath.Base(deployFile.Path), err)
			}
		}
	}

	return nil
}

// 上传节点文件
func (this *Upgrader) uploadNodeDeployFile(ctx context.Context, rpcClient *rpc.RPCClient, path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()

	var buf = make([]byte, 128*4096)
	var isFirst = true

	var hash = md5.New()

	for {
		n, err := fp.Read(buf)
		if n > 0 {
			hash.Write(buf[:n])

			_, uploadErr := rpcClient.APINodeRPC().UploadDeployFileToAPINode(ctx, &pb.UploadDeployFileToAPINodeRequest{
				Filename:     filepath.Base(path),
				Sum:          "",
				ChunkData:    buf[:n],
				IsFirstChunk: isFirst,
				IsLastChunk:  false,
			})
			if uploadErr != nil {
				return uploadErr
			}
			isFirst = false
		}
		if err != nil {
			if err == io.EOF {
				err = nil

				_, uploadErr := rpcClient.APINodeRPC().UploadDeployFileToAPINode(ctx, &pb.UploadDeployFileToAPINodeRequest{
					Filename:     filepath.Base(path),
					Sum:          fmt.Sprintf("%x", hash.Sum(nil)),
					ChunkData:    nil,
					IsFirstChunk: false,
					IsLastChunk:  true,
				})
				if uploadErr != nil {
					return uploadErr
				}

				break
			}
			return err
		}
	}

	return nil
}
