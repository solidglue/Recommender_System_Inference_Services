package model

import (
	"context"
	"infer-microservices/common"
	"infer-microservices/common/flags"
	framework "infer-microservices/common/tensorflow_gogofaster/core/framework"
	pb "infer-microservices/common/tfserving_gogofaster"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"time"

	types "github.com/gogo/protobuf/types"
)

var tfservingModelVersion int64
var tfservingTimeout int64
var modelClientInstance *ModelClient

type ModelClient struct {
	modelName          string                 //model name.
	tfservingModelName string                 //model name of tfserving config list.
	tfservingGrpcPool  *common.GRPCPool       //tfserving grpc pool.
	fieldsSpec         map[string]interface{} //feaure engine conf.
	userRedisKeyPre    string                 //user feature redis key pre.
	itemRedisKeyPre    string                 //item feature redis key pre.
}

func init() {
	flagFactory := flags.FlagFactory{}
	flagTensorflow := flagFactory.FlagTensorflowFactory()

	tfservingModelVersion = *flagTensorflow.GetTfservingModelVersion()
	tfservingTimeout = *flagTensorflow.GetTfservingTimeoutMs()
}

// modelName
func (f *ModelClient) setModelName(modelName string) {
	f.modelName = modelName
}

func (f *ModelClient) GetModelName() string {
	return f.modelName
}

// tfservingModelName
func (f *ModelClient) setTfservingModelName(tfservingModelName string) {
	f.tfservingModelName = tfservingModelName
}

func (f *ModelClient) GetTfservingModelName() string {
	return f.tfservingModelName
}

// tfservingGrpcPool
func (f *ModelClient) setTfservingGrpcPool(tfservingGrpcPool *common.GRPCPool) {
	f.tfservingGrpcPool = tfservingGrpcPool
}

func (f *ModelClient) GetTfservingGrpcPool() *common.GRPCPool {
	return f.tfservingGrpcPool
}

// userRedisKeyPre
func (f *ModelClient) setUserRedisKeyPre(userRedisKeyPre string) {
	f.userRedisKeyPre = userRedisKeyPre
}

func (f *ModelClient) GetUserRedisKeyPre() string {
	return f.userRedisKeyPre
}

// itemRedisKeyPre
func (f *ModelClient) setItemRedisKeyPre(itemRedisKeyPre string) {
	f.itemRedisKeyPre = itemRedisKeyPre
}

func (f *ModelClient) GetItemRedisKeyPre() string {
	return f.itemRedisKeyPre
}

// model conf load
func (m *ModelClient) ConfigLoad(domain string, dataId string, modelConfStr string) error {

	dataConf := utils.Json2Map(modelConfStr)
	for tmpModelName_, tmpModelConf_ := range dataConf { // only 1 model

		modelConfTmp := tmpModelConf_.(map[string]interface{})
		tfservingGrpcConf := modelConfTmp["tfservingGrpcAddr"].(map[string]interface{})
		modelName := tfservingGrpcConf["tfservingModelName"].(string) //tfserving config list modelname

		// create tfserving grpc pool
		tfservingGrpcPool, err := common.CreateGrpcConn(tfservingGrpcConf)
		if err != nil {
			return err
		}

		//fieldsSpec := modelConfTmp["fieldsSpec"].(map[string]interface{})
		userRedisKeyPre := modelConfTmp["userRedisKeyPre"].(string)
		itemRedisKeyPre := modelConfTmp["itemRedisKeyPre"].(string)

		//set
		m.setModelName(tmpModelName_)
		m.setTfservingModelName(modelName)
		m.setTfservingGrpcPool(tfservingGrpcPool)
		m.setUserRedisKeyPre(userRedisKeyPre)
		m.setItemRedisKeyPre(itemRedisKeyPre)
	}

	return nil
}

func (m *ModelClient) requestTfservering(userExamples *[][]byte, userContextExamples *[][]byte, itemExamples *[][]byte, tensorName string) (*[]float32, error) {
	grpcConn, err := m.GetTfservingGrpcPool().Get()
	defer m.GetTfservingGrpcPool().Put(grpcConn)

	if err != nil {
		return nil, err
	}
	predictClient := pb.NewPredictionServiceClient(grpcConn)

	version := &types.Int64Value{Value: tfservingModelVersion}
	predictRequest := &pb.PredictRequest{
		ModelSpec: &pb.ModelSpec{
			Name:    m.GetModelName(),
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
	predictOut, _ := predict.Outputs[tensorName]

	return &predictOut.FloatVal, nil
}

// request embedding vector from tfserving
func (m *ModelClient) Embedding(examples common.ExampleFeatures, tensorName string) (*[]float32, error) {

	userExamples := make([][]byte, 0)
	userContextExamples := make([][]byte, 0)
	itemExamples := make([][]byte, 0)

	userExamples = append(userExamples, *(examples.UserExampleFeatures.Buff))
	userContextExamples = append(userContextExamples, *(examples.UserContextExampleFeatures.Buff))

	response, err := m.requestTfservering(&userExamples, &itemExamples, &userContextExamples, tensorName)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return response, nil
}

// request rank scores from tfserving
func (m *ModelClient) RankPredict(examples common.ExampleFeatures, tensorName string) (*[]string, *[]float32, error) {

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

	scores, err := m.requestTfservering(&userExamples, &userContextExamples, &itemExamples, tensorName)

	if err != nil {
		logs.Error(err)
		return nil, nil, err
	}

	return &items, scores, nil
}
