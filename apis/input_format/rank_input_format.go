package input_format

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/cores/service_config"
)

type RankInputFormat struct {
}

func (d *RankInputFormat) InputCheckAndFormat(in *apis.RecRequest, serverConn *service_config.ServiceConfig) (cores.DeepFM, error) {

	deepfm := cores.DeepFM{}

	//dataid
	if in.GetDataId() == "" {
		return deepfm, errors.New("dataid can not be empty")
	}

	//groupid
	userId := in.GetUserId()
	if userId == "" {
		return deepfm, errors.New("groupid can not be empty")
	}

	//itemlist
	itemList := make([]string, 0)
	itemListIn := in.GetItemList()
	if in.GetItemList() == nil {
		return deepfm, errors.New("itemlist can not be empty")
	} else {
		itemList = itemListIn
	}

	deepfm = cores.DeepFM{}
	deepfm.SetUserId(userId)
	deepfm.SetItemList(itemList)
	deepfm.SetServiceConfig(serverConn)

	return deepfm, nil
}
