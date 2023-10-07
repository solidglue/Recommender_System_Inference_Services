package server

import (
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/utils"
	"sync"

	_ "dubbo.apache.org/dubbo-go/v3/imports" // dubbogo 框架依赖，所有dubbogo进程都需要隐式引入一次
)

//TODO: 传来的参数不固定，且枚举太多，考虑反射(性能差，慎用) https://blog.csdn.net/DkSakura/article/details/116588382

// TODO:此处可以改用策略模式  https://zhuanlan.zhihu.com/p/392843448
var rankWg sync.WaitGroup

type rankServer struct {
	deepfm cores.DeepFM
}

func (r *rankServer) dubboInferServer() (*apis.RecResponse, error) {

	response := apis.RecResponse{}
	response.SetCode(404)

	//关闭go2sky, 此版本走老em召回
	result, err := r.deepfm.RankInferNoSkywalking(nil)
	if err != nil {
		return &response, err
	}

	//召回结果封装
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	rankCh := make(chan string, len(resultList))
	if len(resultList) > 0 {
		//取结果
		for i := 0; i < len(resultList); i++ {
			rankWg.Add(1)
			go fmtRankResponse(resultList[i], rankCh)
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

func fmtRankResponse(itemScore map[string]interface{}, rankCh chan string) {

	defer rankWg.Done()

	itemId := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))

	itemInfo := apis.ItemInfo{}
	itemInfo.SetItemId(itemId)
	itemInfo.SetScore(score)

	itemScoreStr := utils.Struct2Json(itemInfo)
	rankCh <- itemScoreStr

}
