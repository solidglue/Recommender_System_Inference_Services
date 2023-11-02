package common

import (
	"infer-microservices/common/flags"
	"infer-microservices/utils/logs"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func init() {
	flagFactory := flags.FlagFactory{}
	flagViper := flagFactory.CreateFlagViper()

	viper.SetConfigName(*flagViper.GetConfigName())
	viper.SetConfigType(*flagViper.GetConfigType())
	viper.AddConfigPath(*flagViper.GetConfigPath())
}

func loadViperConfigFile() {
	err := viper.ReadInConfig()
	if err != nil {
		logs.Error(err)
	}
	userId := viper.GetString("useId")
	itemId := viper.GetString("itemId")

	bloomPush(GetUserBloomFilterInstance(), userId)
	bloomPush(GetUserBloomFilterInstance(), itemId)
}

func WatchBloomConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logs.Info("Config file changed:", e.Name) //the file only contains new users and new items in past 1 hour.

		loadViperConfigFile()

	})
}
