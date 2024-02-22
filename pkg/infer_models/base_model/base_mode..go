package base_model

import (
	"context"

	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/internal/flags"
	"infer-microservices/internal/logs"
	framework "infer-microservices/internal/tensorflow_gogofaster/core/framework"
	tfserving "infer-microservices/internal/tfserving_gogofaster"
	"infer-microservices/internal/utils"
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/config_loader/model_config"
	"infer-microservices/pkg/infer_samples/feature"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gogo/protobuf/types"
)

var tfservingModelVersion int64
var tfservingTimeout int64
var baseModelInstance *BaseModel
var bigCacheConf bigcache.Config
var bigCacheRsp *bigcache.BigCache
var lifeWindowS time.Duration
var cleanWindowS time.Duration
var hardMaxCacheSize int
var maxEntrySize int
var maxEntriesInWindow int
var verbose bool
var shards int
var err error

type BaseModel struct {
	modelName     string
	serviceConfig *config_loader.ServiceConfig
	tensorName    string
	bigCacheRsp   *bigcache.BigCache
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagTensorflow := flagFactory.CreateFlagTensorflow()
	tfservingModelVersion = *flagTensorflow.GetTfservingModelVersion()
	tfservingTimeout = *flagTensorflow.GetTfservingTimeoutMs()

	flagCache := flagFactory.CreateFlagCache()
	lifeWindowS = time.Duration(*flagCache.GetBigcacheLifeWindowS())
	cleanWindowS = time.Duration(*flagCache.GetBigcacheCleanWindowS())
	hardMaxCacheSize = *flagCache.GetBigcacheHardMaxCacheSize()
	maxEntrySize = *flagCache.GetBigcacheMaxEntrySize()
	bigCacheConf = bigcache.Config{
		Shards:             shards,
		LifeWindow:         lifeWindowS * time.Minute,
		CleanWindow:        cleanWindowS * time.Minute,
		MaxEntriesInWindow: maxEntriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            verbose,
		HardMaxCacheSize:   hardMaxCacheSize,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	//set cache
	bigCacheRsp, err = bigcache.NewBigCache(bigCacheConf)
	if err != nil {
		logs.Error(err)
	}
}

// singleton instance
func init() {
	baseModelInstance = new(BaseModel)
}

func GetBaseModelInstance() *BaseModel {
	return baseModelInstance
}

// modelName
func (b *BaseModel) SetModelName(modelName string) {
	b.modelName = modelName
}

func (b *BaseModel) GetModelName() string {
	return b.modelName
}

// serviceConfig *service_config.ServiceConfig
func (b *BaseModel) SetServiceConfig(serviceConfig *config_loader.ServiceConfig) {
	b.serviceConfig = serviceConfig
}

func (b *BaseModel) GetServiceConfig() *config_loader.ServiceConfig {
	return b.serviceConfig
}

// tensorName
func (b *BaseModel) SetTensorName(tensorName string) {
	b.tensorName = tensorName
}

func (b *BaseModel) GetTensorName() string {
	return b.tensorName
}

// bigCacheRsp
func (b *BaseModel) SetBigCacheRsp(bigCacheRsp *bigcache.BigCache) {
	b.bigCacheRsp = bigCacheRsp
}

func (b *BaseModel) GetBigCacheRsp() *bigcache.BigCache {
	return b.bigCacheRsp
}

// request rank scores from tfserving
func (b *BaseModel) RankPredict(model model_config.ModelConfig, examples feature.ExampleFeatures, tensorName string) (*[]string, *[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)
	items := make([]string, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	for _, itemExample := range *examples.ItemSeqExampleFeatures {
		items = append(items, *(itemExample.Key))
		itemExamples = append(itemExamples, *(itemExample.Buff))
	}
	scores, err := b.requestTfservering(model, &userExamples, &userContextExamples, &itemExamples, tensorName)

	if err != nil {
		return nil, nil, err
	}

	return &items, scores, nil
}

// request embedding vector from tfserving
func (b *BaseModel) Embedding(model model_config.ModelConfig, examples feature.ExampleFeatures, tensorName string) (*[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	response, err := b.requestTfservering(model, &userExamples, &itemExamples, &userContextExamples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return response, nil
}

// request tfserving service by grpc
func (b *BaseModel) requestTfservering(model model_config.ModelConfig, userExamples *[][]byte, userContextExamples *[][]byte, itemExamples *[][]byte, tensorName string) (*[]float32, error) {
	grpcConn, err := model.GetTfservingGrpcPool().Get()
	defer model.GetTfservingGrpcPool().Put(grpcConn)
	if err != nil {
		return nil, err
	}
	predictClient := tfserving.NewPredictionServiceClient(grpcConn)
	version := &types.Int64Value{Value: tfservingModelVersion}
	predictRequest := &tfserving.PredictRequest{
		ModelSpec: &tfserving.ModelSpec{
			Name:    model.GetModelName(),
			Version: version,
		},
		Inputs: make(map[string]*framework.TensorProto),
	}

	//user examples
	tensorProtoUser := &framework.TensorProto{
		Dtype: framework.DataType_DT_STRING,
	}
	tensorProtoUser.TensorShape = &framework.TensorShapeProto{
		Dim: []*framework.TensorShapeProto_Dim{
			{
				Size_: int64(len(*userExamples)),
				Name:  "",
			},
		},
	}
	tensorProtoUser.StringVal = *userExamples
	predictRequest.Inputs["userExamples"] = tensorProtoUser

	//context examples, realtime
	tensorProtoUserContext := &framework.TensorProto{
		Dtype: framework.DataType_DT_STRING,
	}
	tensorProtoUserContext.TensorShape = &framework.TensorShapeProto{
		Dim: []*framework.TensorShapeProto_Dim{
			{
				Size_: int64(len(*userContextExamples)),
				Name:  "",
			},
		},
	}
	tensorProtoUserContext.StringVal = *userContextExamples
	predictRequest.Inputs["userContextExamples"] = tensorProtoUserContext

	//item examples
	tensorProtoItem := &framework.TensorProto{
		Dtype: framework.DataType_DT_STRING,
	}
	tensorProtoItem.TensorShape = &framework.TensorShapeProto{
		Dim: []*framework.TensorShapeProto_Dim{
			{
				Size_: int64(len(*itemExamples)),
				Name:  "",
			},
		},
	}
	tensorProtoItem.StringVal = *itemExamples
	predictRequest.Inputs["itemExamples"] = tensorProtoItem

	predictRequest.OutputFilter = []string{tensorName}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tfservingTimeout)*time.Millisecond)
	defer cancel()

	predict, err := predictClient.Predict(ctx, predictRequest)
	if err != nil {
		return nil, err
	}
	predictOut := predict.Outputs[tensorName]

	return &predictOut.FloatVal, nil
}

func (b *BaseModel) InferResultFormat(recallResult *[]*faiss_index.ItemInfo) (*[]map[string]interface{}, error) {
	recall := make([]map[string]interface{}, 0)
	resultCh := make(chan map[string]interface{}, len(*recallResult))

	for idx := 0; idx < len(*recallResult); idx++ {
		rawCell := (*recallResult)[idx]
		go func(raw_cell_ *faiss_index.ItemInfo) {
			returnCell := make(map[string]interface{})
			returnCell["itemid"] = raw_cell_.ItemId
			returnCell["score"] = utils.FloatRound(raw_cell_.Score, 4)
			resultCh <- returnCell
		}(rawCell)
	}

loop:
	for {
		select {
		case <-time.After(time.Millisecond * 100):
			break loop
		case result := <-resultCh:
			recall = append(recall, result)
		}
	}
	close(resultCh)

	return &recall, nil
}
