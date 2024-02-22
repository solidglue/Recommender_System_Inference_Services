package faiss

import (
	"context"
	"infer-microservices/pkg/infer_samples/feature"

	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/internal/flags"
	"infer-microservices/pkg/config_loader/faiss_config"
	"time"
)

var grpcTimeout int64

func init() {
	flagFactory := flags.FlagFactory{}
	flagTensorflow := flagFactory.CreateFlagTensorflow()
	grpcTimeout = *flagTensorflow.GetTfservingTimeoutMs()
}
func FaissVectorSearch(f *faiss_config.FaissIndexConfig, example feature.ExampleFeatures, vector []float32) ([]*faiss_index.ItemInfo, error) {

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
