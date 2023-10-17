package model

import "infer-microservices/common"

type requestTfserveringInterface interface {
	//get infer samples.
	getInferExampleFeatures() (common.ExampleFeatures, error)

	//request tfserving model.
	requestTfservering(userExamples *[][]byte, userContextExamples *[][]byte, itemExamples *[][]byte, tensorName string) (*[]float32, error)
}
