package apps

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/sizes"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type LogWriter struct {
	fp *os.File
	c  chan string
}

func (this *LogWriter) Init() {
	// 创建目录
	var dir = files.NewFile(Tea.LogDir())
	if !dir.Exists() {
		err := dir.Mkdir()
		if err != nil {
			log.Println("[LOG]create log dir failed: " + err.Error())
		}
	}

	// 打开要写入的日志文件
	var logPath = Tea.LogFile("run.log")
	fp, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("[LOG]open log file failed: " + err.Error())
	} else {
		this.fp = fp
	}

	this.c = make(chan string, 1024)

	// 异步写入文件
	var maxFileSize = 2 * sizes.G // 文件最大尺寸，超出此尺寸则清空
	if fp != nil {
		goman.New(func() {
			var totalSize int64 = 0
			stat, err := fp.Stat()
			if err == nil {
				totalSize = stat.Size()
			}

			for message := range this.c {
				totalSize += int64(len(message))
				_, err := fp.WriteString(timeutil.Format("Y/m/d H:i:s ") + message + "\n")
				if err != nil {
					log.Println("[LOG]write log failed: " + err.Error())
				} else {
					// 如果太大则Truncate
					if totalSize > maxFileSize {
						_ = fp.Truncate(0)
						totalSize = 0
					}
				}
			}
		})
	}
}

func (this *LogWriter) Write(message string) {
	backgroundEnv, _ := os.LookupEnv("EdgeBackground")
	if backgroundEnv != "on" {
		// 文件和行号
		var file string
		var line int
		if Tea.IsTesting() {
			var callDepth = 3
			var ok bool
			_, file, line, ok = runtime.Caller(callDepth)
			if ok {
				file = this.packagePath(file)
			}
		}

		if len(file) > 0 {
			log.Println(message + " (" + file + ":" + strconv.Itoa(line) + ")")
		} else {
			log.Println(message)
		}
	}

	this.c <- message
}

func (this *LogWriter) Close() {
	if this.fp != nil {
		_ = this.fp.Close()
	}

	close(this.c)
}

func (this *LogWriter) packagePath(path string) string {
	var pieces = strings.Split(path, "/")
	if len(pieces) >= 2 {
		return strings.Join(pieces[len(pieces)-2:], "/")
	}
	return path
}
