package infer_features

import (
	"infer-microservices/internal"
)

//TODO: add new solution-B：Query and process features during inference, and then generate samples, More flexible.

//get runtime feature from kafka and save to redis
func init() {
	go internal.KafkaConsumer(runtimeFeature)
}

func runtimeFeature(msgKey string, msgValue string) {

	//parse user action json data from kafka ,extra runtime sequence feature.

	//feature engine and save features to redis.

}
