package serverconfigs

type UnixProtocolConfig struct {
	BaseProtocol `yaml:",inline"`
}

func (this *UnixProtocolConfig) Init() error {
	err := this.InitBase()
	if err != nil {
		return err
	}

	return nil
}
