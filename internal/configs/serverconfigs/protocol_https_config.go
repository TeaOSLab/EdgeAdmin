package serverconfigs

type HTTPSProtocolConfig struct {
	BaseProtocol `yaml:",inline"`
}

func (this *HTTPSProtocolConfig) Init() error {
	err := this.InitBase()
	if err != nil {
		return err
	}

	return nil
}
