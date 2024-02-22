package feature

import (
	"infer-microservices/internal"
)

//TODO: get runtime feature from kafka and save to redis

func init() {
	go internal.KafkaConsumer(runtimeFeature)
}

func runtimeFeature(msgKey string, msgValue string) {

	//parse user action json data from kafka ,extra runtime sequence feature.

	//feature engine and save features to redis.

}
