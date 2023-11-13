package model

import (
	"errors"
	"infer-microservices/apis/io"
	"infer-microservices/common"
	"infer-microservices/cores/model/basemodel"
	"infer-microservices/cores/model/deepfm"
	"infer-microservices/cores/model/dssm"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils/logs"
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

	//dataid
	if in.GetDataId() == "" {
		err := errors.New("dataid can not be empty")
		logs.Error(err)
		return modelStrategyMap[modelName]
	}

	//userid
	userId := in.GetUserId()
	if userId == "" {
		err := errors.New("userid can not be empty")
		logs.Error(err)
		return modelStrategyMap[modelName]
	}

	baseModel = basemodel.GetBaseModelInstance()
	baseModel.SetUserBloomFilter(common.GetUserBloomFilterInstance())
	baseModel.SetItemBloomFilter(common.GetItemBloomFilterInstance())
	baseModel.SetUserId(userId)
	baseModel.SetServiceConfig(serverConn)

	dssmModel := &dssm.Dssm{
		BaseModel: *baseModel,
	}

	modelStrategyMap["dssm"] = dssmModel

	deepfmModel := &deepfm.DeepFM{
		BaseModel: *baseModel,
	}

	//itemlist
	itemList := in.GetItemList()
	if itemList == nil {
		err := errors.New("itemlist can not be empty")
		logs.Error(err)
		return modelStrategyMap[modelName]
	}
	deepfmModel.SetItemList(itemList)

	modelStrategyMap["deepfm"] = deepfmModel

	// m.modelStrategyMap["lr"] = lrModel
	// m.modelStrategyMap["fm"] = fmModel

	return modelStrategyMap[modelName]
}
