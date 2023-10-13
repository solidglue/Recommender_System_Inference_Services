package input_format

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/cores/model/deepfm"
	"infer-microservices/cores/service_config"
)

type RankInputFormat struct {
}

func (d *RankInputFormat) InputCheckAndFormat(in *apis.RecRequest, serverConn *service_config.ServiceConfig) (deepfm.DeepFM, error) {
	deepfmm := deepfm.DeepFM{}

	//dataid
	if in.GetDataId() == "" {
		return deepfmm, errors.New("dataid can not be empty")
	}

	//groupid
	userId := in.GetUserId()
	if userId == "" {
		return deepfmm, errors.New("groupid can not be empty")
	}

	//itemlist
	itemList := make([]string, 0)
	itemListIn := in.GetItemList()
	if itemListIn == nil {
		return deepfmm, errors.New("itemlist can not be empty")
	} else {
		itemList = itemListIn
	}

	deepfmm = deepfm.DeepFM{}
	deepfmm.SetUserId(userId)
	deepfmm.SetItemList(itemList)
	deepfmm.SetServiceConfig(serverConn)

	return deepfmm, nil
}
