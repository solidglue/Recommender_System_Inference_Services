package model

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/cores/model/deepfm"
	"infer-microservices/cores/model/dssm"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils/logs"
	"net/http"
	"strings"
)

var inferModel ModelInferInterface

type ModelInferInterface interface {
	//set params
	SetUserId(userId string)
	SetRetNum(retNum int)
	SetItemList(itemList []string)
	SetServiceConfig(serviceConfig *service_config_loader.ServiceConfig)

	//model infer.
	ModelInferSkywalking(r *http.Request) (map[string]interface{}, error)
	ModelInferNoSkywalking(r *http.Request) (map[string]interface{}, error)
}

type ModelFactory struct {
}

func (m ModelFactory) CreateInferModel(modelName string, in *apis.RecRequest, serverConn *service_config_loader.ServiceConfig) (ModelInferInterface, error) {

	if strings.ToLower(modelName) == "dssm" {
		inferModel = &dssm.Dssm{}
	} else if strings.ToLower(modelName) == "deepfm" {
		inferModel = &deepfm.DeepFM{}
	} else {
		err := errors.New("wrong model")
		logs.Error(err)
		return inferModel, err
	}

	//dataid
	if in.GetDataId() == "" {
		return inferModel, errors.New("dataid can not be empty")
	}

	//userid
	userId := in.GetUserId()
	if userId == "" {
		return inferModel, errors.New("userid can not be empty")
	}

	inferModel.SetUserId(userId)
	inferModel.SetServiceConfig(serverConn)

	//only rank model set itemlist
	if strings.ToLower(modelName) == "deepfm" {
		//itemlist
		itemList := make([]string, 0)
		itemListIn := in.GetItemList()
		if itemListIn == nil {
			return inferModel, errors.New("itemlist can not be empty")
		} else {
			itemList = itemListIn
		}

		inferModel.SetItemList(itemList)
	}

	return inferModel, nil
}
