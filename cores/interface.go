package cores

type ConfigLoadInterface interface {
	ConfigLoad(domain string, dataId string, confStr string) error
}
