package model

import (
	"infer-microservices/internal"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/model/basemodel"
	"infer-microservices/pkg/model/deepfm"
	"infer-microservices/pkg/model/dssm"
	"net/http"
)

var baseModel *basemodel.BaseModel
var modelStrategyMap map[string]ModelStrategyInterface

type ModelStrategyInterface interface {
	//model infer.
	GetModelType() string
	ModelInferSkywalking(requestId string, userId string, itemList []string, r *http.Request, createSample basemodel.CreateSampleCallBackFunc) (map[string]interface{}, error)
	ModelInferNoSkywalking(requestId string, userId string, itemList []string, r *http.Request, createSample basemodel.CreateSampleCallBackFunc) (map[string]interface{}, error)
}

type ModelStrategyFactory struct {
}

func init() {
	modelStrategyMap = make(map[string]ModelStrategyInterface, 0)
}

func SetModelStrategyMap(modelStrategy map[string]ModelStrategyInterface) {
	modelStrategyMap = modelStrategy
}

func GetModelStrategyMap() map[string]ModelStrategyInterface {
	return modelStrategyMap
}

func (m *ModelStrategyFactory) CreateModelStrategy(modelName string, serverConn *config_loader.ServiceConfig) ModelStrategyInterface {
	baseModel = basemodel.GetBaseModelInstance()
	baseModel.SetUserBloomFilter(internal.GetUserBloomFilterInstance())
	baseModel.SetItemBloomFilter(internal.GetItemBloomFilterInstance())
	baseModel.SetServiceConfig(serverConn)

	//dssm model
	dssmModel := &dssm.Dssm{}
	dssmModel.SetBaseModel(*baseModel)
	dssmModel.SetModelType("recall")
	modelStrategyMap["dssm"] = dssmModel

	//deepfm model
	deepfmModel := &deepfm.DeepFM{}
	deepfmModel.SetBaseModel(*baseModel)
	deepfmModel.SetModelType("rank")
	modelStrategyMap["deepfm"] = deepfmModel

	// modelStrategyMap["lr"] = lrModel
	// modelStrategyMap["fm"] = fmModel

	return modelStrategyMap[modelName]
}
