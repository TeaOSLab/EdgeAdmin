package apps

import (
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/gosock/pkg/gosock"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// AppCmd App命令帮助
type AppCmd struct {
	product       string
	version       string
	usage         string
	options       []*CommandHelpOption
	appendStrings []string

	directives []*Directive

	sock *gosock.Sock
}

func NewAppCmd() *AppCmd {
	return &AppCmd{
		sock: gosock.NewTmpSock(teaconst.ProcessName),
	}
}

type CommandHelpOption struct {
	Code        string
	Description string
}

// Product 产品
func (this *AppCmd) Product(product string) *AppCmd {
	this.product = product
	return this
}

// Version 版本
func (this *AppCmd) Version(version string) *AppCmd {
	this.version = version
	return this
}

// Usage 使用方法
func (this *AppCmd) Usage(usage string) *AppCmd {
	this.usage = usage
	return this
}

// Option 选项
func (this *AppCmd) Option(code string, description string) *AppCmd {
	this.options = append(this.options, &CommandHelpOption{
		Code:        code,
		Description: description,
	})
	return this
}

// Append 附加内容
func (this *AppCmd) Append(appendString string) *AppCmd {
	this.appendStrings = append(this.appendStrings, appendString)
	return this
}

// Print 打印
func (this *AppCmd) Print() {
	fmt.Println(this.product + " v" + this.version)

	usage := this.usage
	fmt.Println("Usage:", "\n   "+usage)

	if len(this.options) > 0 {
		fmt.Println("")
		fmt.Println("Options:")

		spaces := 20
		max := 40
		for _, option := range this.options {
			l := len(option.Code)
			if l < max && l > spaces {
				spaces = l + 4
			}
		}

		for _, option := range this.options {
			if len(option.Code) > max {
				fmt.Println("")
				fmt.Println("  " + option.Code)
				option.Code = ""
			}

			fmt.Printf("  %-"+strconv.Itoa(spaces)+"s%s\n", option.Code, ": "+option.Description)
		}
	}

	if len(this.appendStrings) > 0 {
		fmt.Println("")
		for _, s := range this.appendStrings {
			fmt.Println(s)
		}
	}
}

// On 添加指令
func (this *AppCmd) On(arg string, callback func()) {
	this.directives = append(this.directives, &Directive{
		Arg:      arg,
		Callback: callback,
	})
}

// Run 运行
func (this *AppCmd) Run(main func()) {
	// 获取参数
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "-v", "version", "-version", "--version":
			this.runVersion()
			return
		case "?", "help", "-help", "h", "-h":
			this.runHelp()
			return
		case "start":
			this.runStart()
			return
		case "stop":
			this.runStop()
			return
		case "restart":
			this.runRestart()
			return
		case "status":
			this.runStatus()
			return
		}

		// 查找指令
		for _, directive := range this.directives {
			if directive.Arg == args[0] {
				directive.Callback()
				return
			}
		}

		fmt.Println("unknown command '" + args[0] + "'")

		return
	}

	// 日志
	writer := new(LogWriter)
	writer.Init()
	logs.SetWriter(writer)

	// 运行主函数
	main()
}

// 版本号
func (this *AppCmd) runVersion() {
	fmt.Println(this.product+" v"+this.version, "(build: "+runtime.Version(), runtime.GOOS, runtime.GOARCH+")")
}

// 帮助
func (this *AppCmd) runHelp() {
	this.Print()
}

// 启动
func (this *AppCmd) runStart() {
	var pid = this.getPID()
	if pid > 0 {
		fmt.Println(this.product+" already started, pid:", pid)
		return
	}

	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		fmt.Println(this.product+"  start failed:", err.Error())
		return
	}

	fmt.Println(this.product+" started ok, pid:", cmd.Process.Pid)
}

// 停止
func (this *AppCmd) runStop() {
	var pid = this.getPID()
	if pid == 0 {
		fmt.Println(this.product + " not started yet")
		return
	}

	_, _ = this.sock.Send(&gosock.Command{Code: "stop"})

	fmt.Println(this.product+" stopped ok, pid:", types.String(pid))
}

// 重启
func (this *AppCmd) runRestart() {
	this.runStop()
	time.Sleep(1 * time.Second)
	this.runStart()
}

// 状态
func (this *AppCmd) runStatus() {
	var pid = this.getPID()
	if pid == 0 {
		fmt.Println(this.product + " not started yet")
		return
	}

	fmt.Println(this.product + " is running, pid: " + types.String(pid))
}

// 获取当前的PID
func (this *AppCmd) getPID() int {
	if !this.sock.IsListening() {
		return 0
	}

	reply, err := this.sock.Send(&gosock.Command{Code: "pid"})
	if err != nil {
		return 0
	}
	return maps.NewMap(reply.Params).GetInt("pid")
}
