package input_format

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/cores/model/dssm"
	"infer-microservices/cores/service_config"
)

type RecallInputFormat struct {
}

func (d *RecallInputFormat) InputCheckAndFormat(in *apis.RecRequest, serverConn *service_config.ServiceConfig) (dssm.Dssm, error) {
	dssmm := dssm.Dssm{}

	//dataid
	if in.GetDataId() == "" {
		return dssmm, errors.New("dataid can not be empty")
	}

	//UserId
	userId := in.GetUserId()
	if userId == "" {
		return dssmm, errors.New("groupid can not be empty")
	}

	//ret num
	retNum := int(in.GetRecallNum())
	if retNum == 0 {
		return dssmm, errors.New("retNum can not be 0")
	}

	dssmm = dssm.Dssm{}
	dssmm.SetUserId(userId)
	dssmm.SetRetNum(retNum)
	dssmm.SetServiceConfig(serverConn)

	return dssmm, nil
}
