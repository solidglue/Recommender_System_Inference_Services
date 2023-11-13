package faiss_config_loader

import (
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"
	"reflect"
	"strconv"
)

type FaissIndexConfig struct {
	indexName     string                     `validate:"required,unique,min=4,max=10"` //index name.
	faissGrpcPool *utils.GRPCPool            `validate:"required"`                     //faiss  grpc pool.
	faissIndexs   *faiss_index.RecallRequest `validate:"required"`                     // faiss index.
}

// index name
func (f *FaissIndexConfig) setIndexName(indexName string) {
	f.indexName = indexName
}

func (f *FaissIndexConfig) GetIndexName() string {
	return f.indexName
}

// grpc pool
func (f *FaissIndexConfig) setFaissGrpcPool(faissGrpcPool *utils.GRPCPool) {
	f.faissGrpcPool = faissGrpcPool
}

func (f *FaissIndexConfig) GetFaissGrpcPool() *utils.GRPCPool {
	return f.faissGrpcPool
}

// FaissIndexs
func (f *FaissIndexConfig) setFaissIndexs(faissIndexs *faiss_index.RecallRequest) {
	f.faissIndexs = faissIndexs
}

func (f *FaissIndexConfig) GetFaissIndexs() *faiss_index.RecallRequest {
	return f.faissIndexs
}

// @implement ConfigLoadInterface
func (f *FaissIndexConfig) ConfigLoad(dataId string, indexConfStr string) error {
	dataConf := utils.ConvertJsonToStruct(indexConfStr)

	// create faiss grpc pool
	faissGrpcConf := dataConf["faissGrpcAddr"].(map[string]interface{})
	faissGrpcPool, err := utils.CreateGrpcConn(faissGrpcConf)
	if err != nil {
		return err
	}

	//INFO:the recallNum param from http request,maybe int/ float /string。 user reflect to convert to int32.
	recallNum := int32(100)
	indexInfo := dataConf["indexInfo"].(map[string]interface{})
	for indexName, tmpIndexConf := range indexInfo { //only 1 index
		tmpIndexConfMap := tmpIndexConf.(map[string]interface{})

		//recallNum := int32(tmpIndexConfMap["recall_num"].(float64))
		recallNumType := reflect.TypeOf(tmpIndexConfMap["recallNum"])
		recallNumTypeKind := recallNumType.Kind()
		switch recallNumTypeKind {
		case reflect.String:
			recallNumStr, ok := tmpIndexConfMap["recallNum"].(string)
			if ok {
				recallNum64, err := strconv.ParseInt(recallNumStr, 10, 64)
				if err != nil {
					logs.Error(err)
				} else {
					recallNum = int32(recallNum64)
				}
			}
		case reflect.Float32, reflect.Float64, reflect.Int16, reflect.Int, reflect.Int64, reflect.Int8,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			recallNum, _ = tmpIndexConfMap["recallNum"].(int32)
		default:
			logs.Info("unkown type, set recallnum to 100")
		}

		indexName_ := tmpIndexConfMap["index_name"].(string)
		indexInfoStruct := &faiss_index.RecallRequest{
			IndexName: indexName_,
			RecallNum: recallNum,
		}

		f.setIndexName(indexName)
		f.setFaissGrpcPool(faissGrpcPool)
		f.setFaissIndexs(indexInfoStruct)
	}

	return nil
}
