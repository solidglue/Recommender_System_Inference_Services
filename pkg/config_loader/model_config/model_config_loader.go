package model_config

import (
	"infer-microservices/internal"
	"infer-microservices/pkg/utils"
)

type ModelConfig struct {
	modelName          string             `validate:"required,unique,min=4,max=10"` //model name.
	tfservingModelName string             `validate:"required,min=4,max=10"`        //model name of tfserving config list.
	tfservingGrpcPool  *internal.GRPCPool `validate:"required"`                     //tfserving grpc pool.
	//fieldsSpec         map[string]interface{} //feaure engine conf.
	userRedisKeyPre string `validate:"required,min=4,max=10"` //user feature redis key pre.
	itemRedisKeyPre string `validate:"required,min=4,max=10"` //item feature redis key pre.
}

func init() {
}

// modelName
func (f *ModelConfig) setModelName(modelName string) {
	f.modelName = modelName
}

func (f *ModelConfig) GetModelName() string {
	return f.modelName
}

// tfservingModelName
func (f *ModelConfig) setTfservingModelName(tfservingModelName string) {
	f.tfservingModelName = tfservingModelName
}

func (f *ModelConfig) GetTfservingModelName() string {
	return f.tfservingModelName
}

// tfservingGrpcPool
func (f *ModelConfig) setTfservingGrpcPool(tfservingGrpcPool *internal.GRPCPool) {
	f.tfservingGrpcPool = tfservingGrpcPool
}

func (f *ModelConfig) GetTfservingGrpcPool() *internal.GRPCPool {
	return f.tfservingGrpcPool
}

// userRedisKeyPre
func (f *ModelConfig) setUserRedisKeyPre(userRedisKeyPre string) {
	f.userRedisKeyPre = userRedisKeyPre
}

func (f *ModelConfig) GetUserRedisKeyPre() string {
	return f.userRedisKeyPre
}

// itemRedisKeyPre
func (f *ModelConfig) setItemRedisKeyPre(itemRedisKeyPre string) {
	f.itemRedisKeyPre = itemRedisKeyPre
}

func (f *ModelConfig) GetItemRedisKeyPre() string {
	return f.itemRedisKeyPre
}

// @implement ConfigLoadInterface
func (m *ModelConfig) ConfigLoad(dataId string, modelConfStr string) error {

	dataConf := utils.ConvertJsonToStruct(modelConfStr)
	for _, tmpModelConf_ := range dataConf { // only 1 model

		modelConfTmp := tmpModelConf_.(map[string]interface{})
		tfservingGrpcConf := modelConfTmp["tfservingGrpcAddr"].(map[string]interface{})
		modelName := tfservingGrpcConf["tfservingModelName"].(string) //tfserving config list modelname

		// create tfserving grpc pool
		tfservingGrpcPool, err := internal.CreateGrpcConn(tfservingGrpcConf)
		if err != nil {
			return err
		}

		//fieldsSpec := modelConfTmp["fieldsSpec"].(map[string]interface{})
		userRedisKeyPre := modelConfTmp["userRedisKeyPre"].(string)
		itemRedisKeyPre := modelConfTmp["itemRedisKeyPre"].(string)

		//set
		m.setModelName(dataId)
		m.setTfservingModelName(modelName)
		m.setTfservingGrpcPool(tfservingGrpcPool)
		m.setUserRedisKeyPre(userRedisKeyPre)
		m.setItemRedisKeyPre(itemRedisKeyPre)
	}

	return nil
}
