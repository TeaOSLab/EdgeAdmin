package serverconfigs

type UDPProtocolConfig struct {
	BaseProtocol `yaml:",inline"`
}

func (this *UDPProtocolConfig) Init() error {
	err := this.InitBase()
	if err != nil {
		return err
	}

	return nil
}
