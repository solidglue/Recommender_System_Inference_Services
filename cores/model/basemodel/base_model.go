package basemodel

import (
	"context"
	"infer-microservices/common"
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/common/flags"
	framework "infer-microservices/common/tensorflow_gogofaster/core/framework"
	tfserving "infer-microservices/common/tfserving_gogofaster"
	"infer-microservices/cores/service_config_loader"
	"infer-microservices/utils/logs"

	"infer-microservices/utils"
	"sync"
	"time"

	"github.com/gogo/protobuf/types"
)

var tfservingModelVersion int64
var tfservingTimeout int64

type BaseModel struct {
	modelName     string
	userId        string
	serviceConfig *service_config_loader.ServiceConfig
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagTensorflow := flagFactory.CreateFlagTensorflow()
	tfservingModelVersion = *flagTensorflow.GetTfservingModelVersion()
	tfservingTimeout = *flagTensorflow.GetTfservingTimeoutMs()
}

// modelname
// userid
func (b *BaseModel) SetModelName(modelName string) {
	b.modelName = modelName
}

func (b *BaseModel) GetModelName() string {
	return b.modelName
}

// userid
func (b *BaseModel) SetUserId(userId string) {
	b.userId = userId
}

func (b *BaseModel) GetUserId() string {
	return b.userId
}

// serviceConfig *service_config.ServiceConfig
func (b *BaseModel) SetServiceConfig(serviceConfig *service_config_loader.ServiceConfig) {
	b.serviceConfig = serviceConfig
}

func (b *BaseModel) GetServiceConfig() *service_config_loader.ServiceConfig {
	return b.serviceConfig
}

func (d *BaseModel) GetInferExampleFeatures() (common.ExampleFeatures, error) {
	panic("please overwrite in extend models. ")

}

// get user tfrecords offline samples
func (b *BaseModel) GetUserExampleFeatures() (*common.SeqExampleBuff, error) {

	//TODO: update context features.
	redisKey := b.serviceConfig.GetModelConfig().GetUserRedisKeyPre() + b.userId
	userExampleFeats, err := b.serviceConfig.GetRedisConfig().GetRedisPool().Get(redisKey)

	userSeqExampleBuff := common.SeqExampleBuff{}
	userExampleFeatsBuff := make([]byte, 0)

	if err != nil {
		logs.Error("get item features err", err)
		return &userSeqExampleBuff, err
	} else {
		userExampleFeatsBuff = []byte(userExampleFeats) //.(string)
	}

	//protrait features & realtime features.
	userSeqExampleBuff = common.SeqExampleBuff{
		Key:  &b.userId,
		Buff: &userExampleFeatsBuff,
	}

	return &userSeqExampleBuff, nil
}

// get user tfrecords online samples
func (b *BaseModel) GetUserContextExampleFeatures() (*common.SeqExampleBuff, error) {
	userContextSeqExampleBuff := common.SeqExampleBuff{}
	userContextExampleFeatsBuff := make([]byte, 0)

	//TODO: update context features. only from requst. such as location , time
	//context features.
	userContextSeqExampleBuff = common.SeqExampleBuff{
		Key:  &b.userId,
		Buff: &userContextExampleFeatsBuff,
	}

	return &userContextSeqExampleBuff, nil
}

// request tfserving service by grpc
func (b *BaseModel) RequestTfservering(userExamples *[][]byte, userContextExamples *[][]byte, itemExamples *[][]byte, tensorName string) (*[]float32, error) {
	grpcConn, err := b.serviceConfig.GetModelConfig().GetTfservingGrpcPool().Get()
	defer b.serviceConfig.GetModelConfig().GetTfservingGrpcPool().Put(grpcConn)

	if err != nil {
		return nil, err
	}
	predictConfig := tfserving.NewPredictionServiceClient(grpcConn)

	version := &types.Int64Value{Value: tfservingModelVersion}
	predictRequest := &tfserving.PredictRequest{
		ModelSpec: &tfserving.ModelSpec{
			Name:    b.serviceConfig.GetModelConfig().GetModelName(),
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

	predict, err := predictConfig.Predict(ctx, predictRequest)
	if err != nil {
		return nil, err
	}
	predictOut, _ := predict.Outputs[tensorName]

	return &predictOut.FloatVal, nil
}

func (b *BaseModel) InferResultFormat(recallResult *[]*faiss_index.ItemInfo) (*[]map[string]interface{}, error) {
	recall := make([]map[string]interface{}, 0)
	recallTmp := make(chan map[string]interface{}, len(*recallResult)) // 20221011
	var wg sync.WaitGroup

	for idx := 0; idx < len(*recallResult); idx++ {
		rawCell := (*recallResult)[idx]
		wg.Add(1)
		go func(raw_cell_ *faiss_index.ItemInfo) {
			defer wg.Done()
			returnCell := make(map[string]interface{})
			returnCell["itemid"] = raw_cell_.ItemId
			returnCell["score"] = utils.FloatRound(raw_cell_.Score, 4)
			recallTmp <- returnCell
		}(rawCell)
	}
	wg.Wait()
	for idx := 0; idx < len(*recallResult); idx++ {
		returnCellTmp := <-recallTmp
		recall = append(recall, returnCellTmp)
	}
	close(recallTmp)

	return &recall, nil
}
