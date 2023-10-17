package faiss

import (
	"context"
	"infer-microservices/common"
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/common/flags"
	"infer-microservices/cores/service_config_loader/faiss_config_loader"
	"time"
)

var grpcTimeout int64

func init() {
	flagFactory := flags.FlagFactory{}
	flagTensorflow := flagFactory.FlagTensorflowFactory()
	grpcTimeout = *flagTensorflow.GetTfservingTimeoutMs()
}
func FaissVectorSearch(f *faiss_config_loader.FaissIndexClient, example common.ExampleFeatures, vector []float32) ([]*faiss_index.ItemInfo, error) {
	faissIndexs := f.GetFaissIndexs()
	faissGrpcConn, err := f.GetFaissGrpcPool().Get()
	if err != nil {
		return nil, err
	}

	defer f.GetFaissGrpcPool().Put(faissGrpcConn)

	faissClient := faiss_index.NewGrpcRecallServerServiceClient(faissGrpcConn)
	vector_info := faiss_index.UserVectorInfo{
		UserVector: vector,
	}

	index_conf_tmp := &faiss_index.RecallRequest{
		IndexName:       faissIndexs.IndexName,
		UserVectorInfo_: &vector_info,
		RecallNum:       faissIndexs.RecallNum,
	}

	if len(*example.UserExampleFeatures.Buff) == 0 || len(*example.UserContextExampleFeatures.Buff) == 0 {
		return make([]*faiss_index.ItemInfo, 0), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(grpcTimeout)*time.Millisecond)
	defer cancel()

	rst, err := faissClient.GrpcRecall(ctx, index_conf_tmp)
	if err != nil {
		return nil, err
	}

	return rst.ItemInfo_, nil
}
