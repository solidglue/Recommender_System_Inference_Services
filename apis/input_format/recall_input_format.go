package input_format

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/cores/service_config"
)

type RecallInputFormat struct {
}

func (d *RecallInputFormat) InputCheckAndFormat(in *apis.RecRequest, serverConn *service_config.ServiceConfig) (cores.Dssm, error) {

	dssm := cores.Dssm{}

	//dataid
	if in.GetDataId() == "" {
		return dssm, errors.New("dataid can not be empty")
	}

	//UserId
	userId := in.GetUserId()
	if userId == "" {
		return dssm, errors.New("groupid can not be empty")
	}

	//ret num
	retNum := int(in.GetRecallNum())
	if retNum == 0 {
		return dssm, errors.New("retNum can not be 0")
	}

	dssm = cores.Dssm{}
	dssm.SetUserId(userId)
	dssm.SetRetNum(retNum)
	dssm.SetServiceConfig(serverConn)

	return dssm, nil
}
