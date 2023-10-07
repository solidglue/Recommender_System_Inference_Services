package server

import (
	"infer-microservices/apis"
	"infer-microservices/cores"
	"infer-microservices/utils"
	"sync"

	_ "dubbo.apache.org/dubbo-go/v3/imports" // dubbogo 框架依赖，所有dubbogo进程都需要隐式引入一次
)

//TODO: 传来的参数不固定，且枚举太多，考虑反射(性能差，慎用) https://blog.csdn.net/DkSakura/article/details/116588382

//TODO:此处可以改用策略模式  https://zhuanlan.zhihu.com/p/392843448

var recallWg sync.WaitGroup

type recallServer struct {
	dssm cores.Dssm
}

func (d *recallServer) dubboInferServer() (*apis.RecResponse, error) {

	response := apis.RecResponse{}
	response.SetCode(404)

	//关闭go2sky, 此版本走老em召回
	result, err := d.dssm.RecallInferNoSkywalking(nil)
	if err != nil {
		return &response, err
	}

	//TODO:go并发

	//召回结果封装
	itemsScores := make([]string, 0)
	resultList := result["data"].([]map[string]interface{})
	recallCh := make(chan string, len(resultList))

	if len(resultList) > 0 {

		//取结果
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
