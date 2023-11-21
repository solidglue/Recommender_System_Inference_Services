package faiss_config

import (
	"infer-microservices/internal"
	faiss_index "infer-microservices/internal/faiss_gogofaster"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/utils"
	"reflect"
	"strconv"
)

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
func (f *FaissIndexConfig) ConfigLoad(dataId string, indexConfStr string) error {
	dataConf := utils.ConvertJsonToStruct(indexConfStr)

	// create faiss grpc pool
	faissGrpcConf := dataConf["faissGrpcAddr"].(map[string]interface{})
	faissGrpcPool, err := internal.CreateGrpcConn(faissGrpcConf)
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
