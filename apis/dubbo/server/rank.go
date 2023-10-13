package server

import (
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/utils"
	"sync"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

var rankWg sync.WaitGroup

type rankServer struct {
	deepfm cores.DeepFM
}

func (r *rankServer) dubboInferServer() (*apis.RecResponse, error) {
	response := apis.RecResponse{}
	response.SetCode(404)

	//close go2sky in dubbo service .
	//TODO: get go2sky config from config file.
	result, err := r.deepfm.RankInferNoSkywalking(nil)
	if err != nil {
		return &response, err
	}

	//package infer result.
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	rankCh := make(chan string, len(resultList))
	if len(resultList) > 0 {
		for i := 0; i < len(resultList); i++ {
			rankWg.Add(1)
			go formatRankResponse(resultList[i], rankCh)
		}

		rankWg.Wait()
		close(rankCh)
		for itemScore := range rankCh {
			itemsScores = append(itemsScores, itemScore)
		}

		response.SetCode(200)
		response.SetMessage("success")
		response.SetData(itemsScores)
	}

	return &response, err
}

func formatRankResponse(itemScore map[string]interface{}, rankCh chan string) {
	defer rankWg.Done()

	itemId := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))

	itemInfo := apis.ItemInfo{}
	itemInfo.SetItemId(itemId)
	itemInfo.SetScore(score)

	itemScoreStr := utils.Struct2Json(itemInfo)
	rankCh <- itemScoreStr
}
