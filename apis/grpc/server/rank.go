package server

//package main

import (
	grpc_api "infer-microservices/apis/grpc/server/api_gogofaster"
	"infer-microservices/cores"
	"infer-microservices/utils/logs"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// INFO: GO goroutine recall and rank

var rankWg sync.WaitGroup

type rankServer struct {
	deepfm cores.DeepFM
}

func (r *rankServer) grpcInferServer() (*grpc_api.RecommendResponse, error) {
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
				response, err = r.deepfm.RankInferSkywalking(nil)
			} else {
				response, err = r.deepfm.RankInferNoSkywalking(nil)
			}

			if err != nil {
				logs.Error("request tfserving fail:", resp_info)
				return resp_info, err
			}

			result := make([]*grpc_api.ItemInfo, 0)
			resultList := response["data"].([]map[string]interface{})
			rankCh := make(chan *grpc_api.ItemInfo, len(resultList))
			for i := 0; i < len(resultList); i++ {
				rankWg.Add(1)
				go formatRankResponse(resultList[i], rankCh)
			}

			rankWg.Wait()
			close(rankCh)
			for itemScore := range rankCh {
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

func formatRankResponse(itemScore map[string]interface{}, rankCh chan *grpc_api.ItemInfo) {
	defer rankWg.Done()

	itemid := itemScore["itemid"].(string)
	score := float32(itemScore["score"].(float64))
	itemInfo := &grpc_api.ItemInfo{
		Itemid: itemid,
		Score:  score,
	}

	rankCh <- itemInfo
}
