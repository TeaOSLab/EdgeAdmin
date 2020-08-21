package serverconfigs

type TLSProtocolConfig struct {
	BaseProtocol `yaml:",inline"`
}

func (this *TLSProtocolConfig) Init() error {
	err := this.InitBase()
	if err != nil {
		return err
	}

	return nil
}
