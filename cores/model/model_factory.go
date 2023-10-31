package model

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/common"
	"infer-microservices/cores/model/basemodel"
	"infer-microservices/cores/model/deepfm"
	"infer-microservices/cores/model/dssm"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils/logs"
	"net/http"
	"strings"
)

var baseModel *basemodel.BaseModel
var inferModel ModelInferInterface

type ModelInferInterface interface {
	//get infer samples.
	GetInferExampleFeatures() (common.ExampleFeatures, error)

	//model infer.
	ModelInferSkywalking(r *http.Request) (map[string]interface{}, error)
	ModelInferNoSkywalking(r *http.Request) (map[string]interface{}, error)
}

type ModelFactory struct {
}

func (m ModelFactory) CreateInferModel(modelName string, in *apis.RecRequest, serverConn *service_config_loader.ServiceConfig) (ModelInferInterface, error) {

	//dataid
	if in.GetDataId() == "" {
		return inferModel, errors.New("dataid can not be empty")
	}

	//userid
	userId := in.GetUserId()
	if userId == "" {
		return inferModel, errors.New("userid can not be empty")
	}

	baseModel = new(basemodel.BaseModel)
	baseModel.SetUserId(userId)
	baseModel.SetServiceConfig(serverConn)

	if strings.ToLower(modelName) == "dssm" {
		inferModel = &dssm.Dssm{
			BaseModel: *baseModel,
		}
	} else if strings.ToLower(modelName) == "deepfm" {
		inferModel_ := &deepfm.DeepFM{
			BaseModel: *baseModel,
		}

		//itemlist
		itemList := in.GetItemList()
		if itemList == nil {
			return &deepfm.DeepFM{}, errors.New("itemlist can not be empty")
		}
		inferModel_.SetItemList(itemList)
		inferModel = inferModel_

	} else {
		err := errors.New("wrong model")
		logs.Error(err)
		return inferModel, err
	}

	return inferModel, nil
}
