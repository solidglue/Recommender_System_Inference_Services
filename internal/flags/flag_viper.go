package flags

type FlagViper struct {
	configName *string
	configType *string
	configPath *string
}

var flagViperInstance *FlagViper

// singleton instance
func init() {
	flagViperInstance = new(FlagViper)
}

func getFlagViperInstance() *FlagViper {
	return flagViperInstance
}

// SetConfigName
func (s *FlagViper) setConfigName(configName *string) {
	s.configName = configName
}

func (s *FlagViper) GetConfigName() *string {
	return s.configName
}

// SetConfigType
func (s *FlagViper) setConfigType(configType *string) {
	s.configType = configType
}

func (s *FlagViper) GetConfigType() *string {
	return s.configType
}

// SetConfigPath
func (s *FlagViper) setConfigPath(configPath *string) {
	s.configPath = configPath
}

func (s *FlagViper) GetConfigPath() *string {
	return s.configPath
}
