package api

import (
	"errors"
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/cores/service_config"
	"infer-microservices/utils"

	_ "dubbo.apache.org/dubbo-go/v3/imports" // dubbogo 框架依赖，所有dubbogo进程都需要隐式引入一次
)

//TODO: 传来的参数不固定，且枚举太多，考虑反射(性能差，慎用) https://blog.csdn.net/DkSakura/article/details/116588382

//TODO:此处可以改用策略模式  https://zhuanlan.zhihu.com/p/392843448

type RankServer struct {
}

func (d *RankServer) inputCheckAndBuild(in *apis.RecRequest, serverConn *service_config.ServiceConfig) (cores.DeepFM, error) {

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

func (r *RankServer) dubboRankInferServer(deepfm cores.DeepFM) (*apis.RecResponse, error) {

	response := apis.RecResponse{}
	response.SetCode(404)

	//关闭go2sky, 此版本走老em召回
	result, err := deepfm.RankInferNoSkywalking(nil)
	if err != nil {
		return &response, err
	}

	//召回结果封装
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			itemId := resultList[i]["itemid"].(string)
			score := float32(resultList[i]["score"].(float64))

			itemInfo := apis.ItemInfo{}
			itemInfo.SetItemId(itemId)
			itemInfo.SetScore(score)

			itemScore := utils.Struct2Json(itemInfo)
			itemsScores = append(itemsScores, itemScore)
		}

		response.SetCode(200)
		response.SetMessage("success")
		response.SetData(itemsScores)
	}

	return &response, err
}
