package service_config_loader

type ConfigLoadInterface interface {
	ConfigLoad(domain string, dataId string, confStr string) error
}
