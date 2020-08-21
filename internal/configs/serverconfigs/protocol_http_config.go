package serverconfigs

type HTTPProtocolConfig struct {
	BaseProtocol `yaml:",inline"`
}

func (this *HTTPProtocolConfig) Init() error {
	err := this.InitBase()
	if err != nil {
		return err
	}

	return nil
}
