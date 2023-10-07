package server

import "infer-microservices/apis"

type dubboInferInterface interface {
	dubboInferServer() (*apis.RecResponse, error)
}
