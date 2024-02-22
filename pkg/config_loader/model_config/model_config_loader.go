package model_config

import (
	"infer-microservices/internal"
	"infer-microservices/internal/utils"
)

// TODO:样本分离，召回、粗排、精排可能用不一样的样本。或者配置哪个用哪些样本，如果每个模型都按照example预存储redis，存储量较大。暂时按照预存redis方案，用空间换时间
type ModelConfig struct {
	modelName          string             `validate:"required,unique,min=4,max=10"` //model name, dssm 、 deepfm
	modelType          string             `validate:"required,unique,min=4,max=10"` //recall 、rank
	tfservingModelName string             `validate:"required,min=4,max=10"`        //model name of tfserving config list.
	tfservingGrpcPool  *internal.GRPCPool `validate:"required"`                     //tfserving grpc pool.
	//fieldsSpec         map[string]interface{} //feaure engine conf.
	userRedisKeyPreOffline  string   `validate:"required,min=4,max=10"` //user offline feature redis key pre.
	userRedisKeyPreRealtime string   `validate:"required,min=4,max=10"` //user Realtime feature redis key pre.
	itemRedisKeyPre         string   `validate:"required,min=4,max=10"` //item feature redis key pre.
	featureList             []string // 召回、粗排、精排等，根据features list选特征
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

// modelType
func (f *ModelConfig) setModelType(modelType string) {
	f.modelType = modelType
}

func (f *ModelConfig) GetModelType() string {
	return f.modelType
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
func (f *ModelConfig) setUserRedisKeyPreOffline(userRedisKeyPreOffline string) {
	f.userRedisKeyPreOffline = userRedisKeyPreOffline
}

func (f *ModelConfig) GetUserRedisKeyPreOffline() string {
	return f.userRedisKeyPreOffline
}

// userRedisKeyPre
func (f *ModelConfig) setUserRedisKeyPreRealtime(userRedisKeyPreRealtime string) {
	f.userRedisKeyPreRealtime = userRedisKeyPreRealtime
}

func (f *ModelConfig) GetUserRedisKeyPreRealtime() string {
	return f.userRedisKeyPreRealtime
}

// itemRedisKeyPre
func (f *ModelConfig) setItemRedisKeyPre(itemRedisKeyPre string) {
	f.itemRedisKeyPre = itemRedisKeyPre
}

func (f *ModelConfig) GetItemRedisKeyPre() string {
	return f.itemRedisKeyPre
}

// featureList
func (f *ModelConfig) setFeatureList(featureList []string) {
	f.featureList = featureList
}

func (f *ModelConfig) GetFeatureList() []string {
	return f.featureList
}

// @implement ConfigLoadInterface
func (m *ModelConfig) ConfigLoad(dataId string, modelConfStr string) error {

	modelConfTmp := utils.ConvertJsonToStruct(modelConfStr)
	tfservingGrpcConf := modelConfTmp["tfservingGrpcAddr"].(map[string]interface{})
	modelName := tfservingGrpcConf["tfservingModelName"].(string) //tfserving config list modelname

	// create tfserving grpc pool
	tfservingGrpcPool, err := internal.CreateGrpcConn(tfservingGrpcConf)
	if err != nil {
		return err
	}

	//fieldsSpec := modelConfTmp["fieldsSpec"].(map[string]interface{})
	userRedisKeyPreOffline := modelConfTmp["userRedisKeyPreOffline"].(string)
	userRedisKeyPreRealtime := modelConfTmp["userRedisKeyPreRealtime"].(string)
	itemRedisKeyPre := modelConfTmp["itemRedisKeyPre"].(string)
	featureList := modelConfTmp["featureList"].([]string)

	//set
	m.setModelName(dataId)
	m.setTfservingModelName(modelName)
	m.setTfservingGrpcPool(tfservingGrpcPool)
	m.setUserRedisKeyPreOffline(userRedisKeyPreOffline)
	m.setUserRedisKeyPreRealtime(userRedisKeyPreRealtime)
	m.setItemRedisKeyPre(itemRedisKeyPre)
	m.setFeatureList(featureList)

	return nil
}
