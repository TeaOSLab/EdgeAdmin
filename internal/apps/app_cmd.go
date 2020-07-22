package apps

import (
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// App命令帮助
type AppCmd struct {
	product       string
	version       string
	usage         string
	options       []*CommandHelpOption
	appendStrings []string

	directives []*Directive
}

func NewAppCmd() *AppCmd {
	return &AppCmd{}
}

type CommandHelpOption struct {
	Code        string
	Description string
}

// 产品
func (this *AppCmd) Product(product string) *AppCmd {
	this.product = product
	return this
}

// 版本
func (this *AppCmd) Version(version string) *AppCmd {
	this.version = version
	return this
}

// 使用方法
func (this *AppCmd) Usage(usage string) *AppCmd {
	this.usage = usage
	return this
}

// 选项
func (this *AppCmd) Option(code string, description string) *AppCmd {
	this.options = append(this.options, &CommandHelpOption{
		Code:        code,
		Description: description,
	})
	return this
}

// 附加内容
func (this *AppCmd) Append(appendString string) *AppCmd {
	this.appendStrings = append(this.appendStrings, appendString)
	return this
}

// 打印
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

// 添加指令
func (this *AppCmd) On(arg string, callback func()) {
	this.directives = append(this.directives, &Directive{
		Arg:      arg,
		Callback: callback,
	})
}

// 运行
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

	// 记录PID
	_ = this.writePid()

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
	proc := this.checkPid()
	if proc != nil {
		fmt.Println(this.product+" already started, pid:", proc.Pid)
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
	proc := this.checkPid()
	if proc == nil {
		fmt.Println(this.product + " not started yet")
		return
	}

	// 停止进程
	_ = proc.Kill()

	// 在Windows上经常不能及时释放资源
	_ = DeletePid(Tea.Root + "/bin/pid")
	fmt.Println(this.product+" stopped ok, pid:", proc.Pid)
}

// 重启
func (this *AppCmd) runRestart() {
	this.runStop()
	time.Sleep(1 * time.Second)
	this.runStart()
}

// 状态
func (this *AppCmd) runStatus() {
	proc := this.checkPid()
	if proc == nil {
		fmt.Println(this.product + " not started yet")
	} else {
		fmt.Println(this.product + " is running, pid: " + fmt.Sprintf("%d", proc.Pid))
	}
}

// 检查PID
func (this *AppCmd) checkPid() *os.Process {
	return CheckPid(Tea.Root + "/bin/pid")
}

// 写入PID
func (this *AppCmd) writePid() error {
	return WritePid(Tea.Root + "/bin/pid")
}
