package infer_samples

import (
	"infer-microservices/internal/flags"
	"infer-microservices/internal/logs"

	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var viperConfig *viper.Viper

// var baseModelObserver Observer
var inferSampleSubject *SampleSubject

func init() {
	flagFactory := flags.FlagFactory{}
	flagViper := flagFactory.CreateFlagViper()

	viperConfig = viper.New()
	viperConfig.SetConfigName(*flagViper.GetConfigName())
	viperConfig.SetConfigType(*flagViper.GetConfigType())
	viperConfig.AddConfigPath(*flagViper.GetConfigPath())

	// //callback func config
	// basemodel0 := BaseModel{}
	// SampleCallBackFuncMap["recall"] = basemodel0.GetInferExampleFeaturesNotContainItems
	// SampleCallBackFuncMap["rank"] = basemodel0.GetInferExampleFeaturesContainItems
}

// TODO:listen kafka,not file.
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
		BloomPush(GetUserBloomFilterInstance(), userId)
	}

	for _, itemId := range itemIdList {
		BloomPush(GetUserBloomFilterInstance(), itemId)
	}

	// //update bloom filter
	inferSampleSubject.NotifyObservers()
}

func WatchBloomConfig() {
	viperConfig.WatchConfig()
	viperConfig.OnConfigChange(func(e fsnotify.Event) {
		logs.Info("Config file changed:", e.Name) //the file only contains new users and new items in past 1 hour.
		loadViperConfigFile()
	})
}
