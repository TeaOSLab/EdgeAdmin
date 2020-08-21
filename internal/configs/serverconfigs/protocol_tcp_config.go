package serverconfigs

type TCPProtocolConfig struct {
	BaseProtocol `yaml:",inline"`
}

func (this *TCPProtocolConfig) Init() error {
	err := this.InitBase()
	if err != nil {
		return err
	}

	return nil
}
