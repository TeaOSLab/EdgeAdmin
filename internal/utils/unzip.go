package utils

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

type Unzip struct {
	zipFile   string
	targetDir string
}

func NewUnzip(zipFile string, targetDir string) *Unzip {
	return &Unzip{
		zipFile:   zipFile,
		targetDir: targetDir,
	}
}

func (this *Unzip) Run() error {
	if len(this.zipFile) == 0 {
		return errors.New("zip file should not be empty")
	}
	if len(this.targetDir) == 0 {
		return errors.New("target dir should not be empty")
	}

	reader, err := zip.OpenReader(this.zipFile)
	if err != nil {
		return err
	}

	defer func() {
		_ = reader.Close()
	}()

	for _, file := range reader.File {
		var info = file.FileInfo()
		var target = this.targetDir + "/" + file.Name

		// 目录
		if info.IsDir() {
			stat, err := os.Stat(target)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				} else {
					err = os.MkdirAll(target, info.Mode())
					if err != nil {
						return err
					}
				}
			} else if !stat.IsDir() {
				err = os.MkdirAll(target, info.Mode())
				if err != nil {
					return err
				}
			}
			continue
		}

		// 文件
		err = func(file *zip.File, target string) error {
			fileReader, err := file.Open()
			if err != nil {
				return err
			}
			defer func() {
				_ = fileReader.Close()
			}()

			// remove old
			_ = os.Remove(target)

			// create new
			fileWriter, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer func() {
				_ = fileWriter.Close()
			}()

			_, err = io.Copy(fileWriter, fileReader)
			return err
		}(file, target)
		if err != nil {
			return err
		}
	}

	return nil
}
