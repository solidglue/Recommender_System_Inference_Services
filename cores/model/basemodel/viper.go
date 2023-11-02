package basemodel

import (
	"infer-microservices/common"
	"infer-microservices/common/flags"
	"infer-microservices/utils/logs"

	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var viperConfig *viper.Viper
var baseModelObserver Observer
var subject *modelSubject

func init() {
	flagFactory := flags.FlagFactory{}
	flagViper := flagFactory.CreateFlagViper()

	viperConfig = viper.New()
	viperConfig.SetConfigName(*flagViper.GetConfigName())
	viperConfig.SetConfigType(*flagViper.GetConfigType())
	viperConfig.AddConfigPath(*flagViper.GetConfigPath())
}

func loadViperConfigFile() {
	err := viperConfig.ReadInConfig()
	if err != nil {
		logs.Error(err)
	}
	userIdStr := viperConfig.GetString("userIdList") // "user1,user2,user3"
	itemIdStr := viperConfig.GetString("itemIdList") //"item1,item2,item3"

	userIdList := strings.Split(userIdStr, ",")
	itemIdList := strings.Split(itemIdStr, ",")

	for _, userId := range userIdList {
		common.BloomPush(common.GetUserBloomFilterInstance(), userId)
	}

	for _, itemId := range itemIdList {
		common.BloomPush(common.GetUserBloomFilterInstance(), itemId)
	}

	// //update bloom filter
	subject = &modelSubject{}
	baseModelObserver = GetBaseModelInstance()
	subject.AddObserver(baseModelObserver)
	subject.NotifyObservers()
}

func WatchBloomConfig() {
	viperConfig.WatchConfig()
	viperConfig.OnConfigChange(func(e fsnotify.Event) {
		logs.Info("Config file changed:", e.Name) //the file only contains new users and new items in past 1 hour.
		loadViperConfigFile()
	})
}
