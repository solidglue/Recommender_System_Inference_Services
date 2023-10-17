package model

import "infer-microservices/common"

type RequestTfserveringInterface interface {
	//request tfserving model.
	RequestTfservering(userExamples *[][]byte, userContextExamples *[][]byte, itemExamples *[][]byte, tensorName string) (*[]float32, error)
}

type GetInferExampleFeaturesInterface interface {
	//get infer samples.
	GetInferExampleFeatures() (common.ExampleFeatures, error)
}
