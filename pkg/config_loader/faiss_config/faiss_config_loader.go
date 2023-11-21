package faiss_config

import (
	"infer-microservices/internal"
	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/utils"
	"reflect"
	"strconv"
)

type FaissIndexConfigs struct {
	faissIndexConfigs []FaissIndexConfig
}

func (f *FaissIndexConfigs) SetFaissIndexConfig(faissIndexConfigs []FaissIndexConfig) {
	f.faissIndexConfigs = faissIndexConfigs
}

func (f *FaissIndexConfigs) GetFaissIndexConfig() []FaissIndexConfig {
	return f.faissIndexConfigs
}

type FaissIndexConfig struct {
	indexName     string                     `validate:"required,unique,min=4,max=10"` //index name.
	faissGrpcPool *internal.GRPCPool         `validate:"required"`                     //faiss  grpc pool.
	faissIndexs   *faiss_index.RecallRequest `validate:"required"`                     // faiss index.
	recallNum     int                        `validate:"required"`                     // faiss recall num.
}

// index name
func (f *FaissIndexConfig) setIndexName(indexName string) {
	f.indexName = indexName
}

func (f *FaissIndexConfig) GetIndexName() string {
	return f.indexName
}

// grpc pool
func (f *FaissIndexConfig) setFaissGrpcPool(faissGrpcPool *internal.GRPCPool) {
	f.faissGrpcPool = faissGrpcPool
}

func (f *FaissIndexConfig) GetFaissGrpcPool() *internal.GRPCPool {
	return f.faissGrpcPool
}

// FaissIndexs
func (f *FaissIndexConfig) setFaissIndexs(faissIndexs *faiss_index.RecallRequest) {
	f.faissIndexs = faissIndexs
}

func (f *FaissIndexConfig) GetFaissIndexs() *faiss_index.RecallRequest {
	return f.faissIndexs
}

// recall num
func (f *FaissIndexConfig) SetRecallNum(recallNum int) {
	f.recallNum = recallNum
}

func (f *FaissIndexConfig) GetRecallNum() int {
	return f.recallNum
}

// @implement ConfigLoadInterface
func (f *FaissIndexConfigs) ConfigLoad(dataId string, indexConfStr string) error {
	dataConf := utils.ConvertJsonToStruct(indexConfStr)

	faissIndexConfigs := make([]FaissIndexConfig, 0)
	// create faiss grpc pool
	faissGrpcConf := dataConf["faissGrpcAddr"].(map[string]interface{})
	faissGrpcPool, err := internal.CreateGrpcConn(faissGrpcConf)
	if err != nil {
		return err
	}

	//INFO:the recallNum param from http request,maybe int/ float /string。 user reflect to convert to int32.
	//INFO:Processing multiple recalls simultaneously to save network overhead
	recallNum := int32(100)
	indexInfo := dataConf["indexInfo"].([]interface{})
	for _, tmpIndexConf := range indexInfo {
		tmpIndexConfMap := tmpIndexConf.(map[string]interface{})
		faissIndexConfig := FaissIndexConfig{}

		//recallNum := int32(tmpIndexConfMap["recall_num"].(float64))
		recallNumType := reflect.TypeOf(tmpIndexConfMap["recallNum"])
		recallNumTypeKind := recallNumType.Kind()
		switch recallNumTypeKind {
		case reflect.String:
			recallNumStr, ok := tmpIndexConfMap["recallNum"].(string)
			if ok {
				recallNum64, err := strconv.ParseInt(recallNumStr, 10, 64)
				if err != nil {
					logs.Warn(err)
				} else {
					recallNum = int32(recallNum64)
				}
			}
		case reflect.Float32, reflect.Float64, reflect.Int16, reflect.Int, reflect.Int64, reflect.Int8,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			recallNum, _ = tmpIndexConfMap["recallNum"].(int32)
		default:
			logs.Warn("unkown type, set recallnum to 100")
		}

		indexName_ := tmpIndexConfMap["indexName"].(string)
		indexInfoStruct := &faiss_index.RecallRequest{
			IndexName: indexName_,
			RecallNum: recallNum,
		}

		faissIndexConfig.setIndexName(indexName_)
		faissIndexConfig.setFaissGrpcPool(faissGrpcPool)
		faissIndexConfig.setFaissIndexs(indexInfoStruct)
		faissIndexConfigs = append(faissIndexConfigs, faissIndexConfig)

		f.SetFaissIndexConfig(faissIndexConfigs)
	}

	return nil
}
