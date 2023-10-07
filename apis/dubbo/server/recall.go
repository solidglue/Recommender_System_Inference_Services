package server

import (
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/utils"
	"sync"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

var recallWg sync.WaitGroup

type recallServer struct {
	dssm cores.Dssm
}

func (d *recallServer) dubboInferServer() (*apis.RecResponse, error) {
	response := apis.RecResponse{}
	response.SetCode(404)

	//close go2sky in dubbo service .
	result, err := d.dssm.RecallInferNoSkywalking(nil)
	if err != nil {
		return &response, err
	}

	//package infer result.
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	recallCh := make(chan string, len(resultList))

	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			recallWg.Add(1)
			go fmtRecallResponse(resultList[i], recallCh)
		}

		recallWg.Wait()
		close(recallCh)
		for itemScore := range recallCh {
			itemsScores = append(itemsScores, itemScore)
		}

		response.SetCode(200)
		response.SetMessage("success")
		response.SetData(itemsScores)
	}

	return &response, err
}

func fmtRecallResponse(itemScore map[string]interface{}, recallCh chan string) {
	defer recallWg.Done()

	itemId := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))

	itemInfo := apis.ItemInfo{}
	itemInfo.SetItemId(itemId)
	itemInfo.SetScore(score)

	itemScoreStr := utils.Struct2Json(itemInfo)
	recallCh <- itemScoreStr
}
