package infer_model

import (
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/config_loader/model_config"
	feature "infer-microservices/pkg/infer_features"
	"infer-microservices/pkg/infer_models/base_model"
	"infer-microservices/pkg/infer_models/ranking/heavy_ranker"
	dssm "infer-microservices/pkg/infer_models/recall/u2i"
	"net/http"
)

var baseModel *base_model.BaseModel
var modelStrategyMap map[string]ModelStrategyInterface

type ModelStrategyInterface interface {
	//model infer.
	GetModelType() string
	ModelInferSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, r *http.Request, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error)
	ModelInferNoSkywalking(model model_config.ModelConfig, requestId string, exposureList []string, inferSample feature.ExampleFeatures, retNum int) (map[string][]map[string]interface{}, error)
}

type ModelStrategyFactory struct {
}

func init() {
	modelStrategyMap = make(map[string]ModelStrategyInterface, 0)
}

func SetModelStrategyMap(modelStrategy map[string]ModelStrategyInterface) {
	modelStrategyMap = modelStrategy
}

func GetModelStrategyMap() map[string]ModelStrategyInterface {
	return modelStrategyMap
}

// TODO: 改成传参行为，便于pipline
func (m *ModelStrategyFactory) CreateModelStrategy(modelName string, serverConn *config_loader.ServiceConfig) ModelStrategyInterface {
	baseModel = base_model.GetBaseModelInstance()
	// baseModel.SetUserBloomFilter(internal.GetUserBloomFilterInstance())
	// baseModel.SetItemBloomFilter(internal.GetItemBloomFilterInstance())
	baseModel.SetServiceConfig(serverConn)

	//dssm model
	dssmModel := &dssm.Dssm{}
	dssmModel.SetBaseModel(*baseModel)
	dssmModel.SetModelType("recall")
	modelStrategyMap["dssm"] = dssmModel
	//	tensorName := "user_embedding"

	//DIN model
	DINModel := &heavy_ranker.DIN{}
	DINModel.SetBaseModel(*baseModel)
	DINModel.SetModelType("rank")
	modelStrategyMap["DIN"] = DINModel
	//	tensorName := "scores"

	// modelStrategyMap["lr"] = lrModel
	// modelStrategyMap["fm"] = fmModel

	return modelStrategyMap[modelName]
}
