package server

import (
	"infer-microservices/cores"
	"sync"
	"time"

	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/utils/logs"

	"golang.org/x/net/context"
)

var recallWg sync.WaitGroup

type recallServer struct {
	dssm cores.Dssm
}

func (r *recallServer) grpcInferServer() (*grpc_api.RecommendResponse, error) {
	resp_info := &grpc_api.RecommendResponse{
		Code: 404,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*150)
	defer cancelFunc()

	for {
		select {
		case <-ctx.Done():
			logs.Info("context timeout.")
			return resp_info, ctx.Err()
		default:
			var response map[string]interface{}
			var err error
			if skywalkingWeatherOpen {
				response, err = r.dssm.RecallInferSkywalking(nil)
			} else {
				response, err = r.dssm.RecallInferNoSkywalking(nil)
			}
			if err != nil {
				logs.Error("request tfserving fail:", resp_info)
				return resp_info, err
			}

			result := make([]*grpc_api.ItemInfo, 0)
			resultList := response["data"].([]map[string]interface{})
			recallCh := make(chan *grpc_api.ItemInfo, len(resultList))
			for i := 0; i < len(resultList); i++ {
				recallWg.Add(1)
				go fmtRecallResponse(resultList[i], recallCh)
			}

			recallWg.Wait()
			close(recallCh)
			for itemScore := range recallCh {
				result = append(result, itemScore)
			}

			resp_info = &grpc_api.RecommendResponse{
				Code:    200,
				Message: "success",
				Data: &grpc_api.ItemInfoList{
					Iteminfo_: result,
				},
			}

			return resp_info, nil
		}
	}
}

func fmtRecallResponse(itemScore map[string]interface{}, rankCh chan *grpc_api.ItemInfo) {

	defer recallWg.Done()

	itemid := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))

	itemInfo := &grpc_api.ItemInfo{
		Itemid: itemid,
		Score:  score,
	}

	rankCh <- itemInfo

}
