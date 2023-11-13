package model

import (
	"infer-microservices/apis/io"
	"infer-microservices/common"
	"infer-microservices/cores/model/basemodel"
	"infer-microservices/cores/model/deepfm"
	"infer-microservices/cores/model/dssm"
	"infer-microservices/cores/service_config_loader"
	"net/http"
)

var baseModel *basemodel.BaseModel
var modelStrategyMap map[string]ModelStrategyInterface

type ModelStrategyInterface interface {
	//model infer.
	ModelInferSkywalking(r *http.Request) (map[string]interface{}, error)
	ModelInferNoSkywalking(r *http.Request) (map[string]interface{}, error)
}

type ModelStrategyFactory struct {
}

func init() {
	modelStrategyMap = make(map[string]ModelStrategyInterface, 0)
}

func (m *ModelStrategyFactory) CreateModelStrategy(modelName string, in *io.RecRequest, serverConn *service_config_loader.ServiceConfig) ModelStrategyInterface {
	baseModel = basemodel.GetBaseModelInstance()
	baseModel.SetUserBloomFilter(common.GetUserBloomFilterInstance())
	baseModel.SetItemBloomFilter(common.GetItemBloomFilterInstance())
	baseModel.SetUserId(in.GetUserId())
	baseModel.SetServiceConfig(serverConn)

	//dssm model
	dssmModel := &dssm.Dssm{
		BaseModel: *baseModel,
	}
	modelStrategyMap["dssm"] = dssmModel

	//deepfm model
	deepfmModel := &deepfm.DeepFM{
		BaseModel: *baseModel,
	}
	deepfmModel.SetItemList(in.GetItemList())
	modelStrategyMap["deepfm"] = deepfmModel

	// modelStrategyMap["lr"] = lrModel
	// modelStrategyMap["fm"] = fmModel

	return modelStrategyMap[modelName]
}
