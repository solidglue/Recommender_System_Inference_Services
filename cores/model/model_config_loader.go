package model

import (
	"infer-microservices/common/flags"
	"infer-microservices/utils"
)

var tfservingModelVersion int64
var tfservingTimeout int64
var modelClientInstance *ModelClient

type ModelClient struct {
	modelName          string                 //model name.
	tfservingModelName string                 //model name of tfserving config list.
	tfservingGrpcPool  *utils.GRPCPool        //tfserving grpc pool.
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
func (f *ModelClient) setTfservingGrpcPool(tfservingGrpcPool *utils.GRPCPool) {
	f.tfservingGrpcPool = tfservingGrpcPool
}

func (f *ModelClient) GetTfservingGrpcPool() *utils.GRPCPool {
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
		tfservingGrpcPool, err := utils.CreateGrpcConn(tfservingGrpcConf)
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
