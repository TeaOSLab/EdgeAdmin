//go:build windows

package utils

import (
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"os/exec"
)

// 安装服务
func (this *ServiceManager) Install(exePath string, args []string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connecting: %w please 'Run as administrator' again", err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(this.Name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", this.Name)
	}

	s, err = m.CreateService(this.Name, exePath, mgr.Config{
		DisplayName: this.Name,
		Description: this.Description,
		StartType:   windows.SERVICE_AUTO_START,
	}, args...)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	defer s.Close()

	return nil
}

// 启动服务
func (this *ServiceManager) Start() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(this.Name)
	if err != nil {
		return fmt.Errorf("could not access service: %w", err)
	}
	defer s.Close()
	err = s.Start("service")
	if err != nil {
		return fmt.Errorf("could not start service: %w", err)
	}

	return nil
}

// 删除服务
func (this *ServiceManager) Uninstall() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connecting: %w please 'Run as administrator' again", err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(this.Name)
	if err != nil {
		return fmt.Errorf("open service: %w", err)
	}

	// shutdown service
	_, err = s.Control(svc.Stop)
	if err != nil {
		fmt.Printf("shutdown service: %s\n", err.Error())
	}

	defer s.Close()
	err = s.Delete()
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}
	return nil
}

// 运行
func (this *ServiceManager) Run() {
	err := svc.Run(this.Name, this)
	if err != nil {
		this.LogError(err.Error())
	}
}

// 同服务管理器的交互
func (this *ServiceManager) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	changes <- svc.Status{
		State: svc.StartPending,
	}

	changes <- svc.Status{
		State:   svc.Running,
		Accepts: cmdsAccepted,
	}

	// start service
	this.Log("start")
	this.cmdStart()

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				this.Log("cmd: Interrogate")
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				this.Log("cmd: Stop|Shutdown")

				// stop service
				this.cmdStop()

				break loop
			case svc.Pause:
				this.Log("cmd: Pause")

				// stop service
				this.cmdStop()

				changes <- svc.Status{
					State:   svc.Paused,
					Accepts: cmdsAccepted,
				}
			case svc.Continue:
				this.Log("cmd: Continue")

				// start service
				this.cmdStart()

				changes <- svc.Status{
					State:   svc.Running,
					Accepts: cmdsAccepted,
				}
			default:
				this.LogError(fmt.Sprintf("unexpected control request #%d\r\n", c))
			}
		}
	}
	changes <- svc.Status{
		State: svc.StopPending,
	}
	return
}

// 启动Web服务
func (this *ServiceManager) cmdStart() {
	cmd := exec.Command(Tea.Root+Tea.DS+"bin"+Tea.DS+teaconst.SystemdServiceName+".exe", "start")
	err := cmd.Start()
	if err != nil {
		this.LogError(err.Error())
	}
}

// 停止Web服务
func (this *ServiceManager) cmdStop() {
	cmd := exec.Command(Tea.Root+Tea.DS+"bin"+Tea.DS+teaconst.SystemdServiceName+".exe", "stop")
	err := cmd.Start()
	if err != nil {
		this.LogError(err.Error())
	}
}
