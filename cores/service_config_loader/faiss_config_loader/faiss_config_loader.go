package faiss_config_loader

import (
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/utils"
)

type FaissIndexConfig struct {
	indexName     string                     //index name.
	faissGrpcPool *utils.GRPCPool            //faiss  grpc pool.
	faissIndexs   *faiss_index.RecallRequest // faiss index.
}

// INFO: singleton instance
func init() {

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

// faiss index conf load
func (f *FaissIndexConfig) ConfigLoad(domain string, dataId string, indexConfStr string) error {
	dataConf := utils.ConvertJsonToStruct(indexConfStr)

	// create faiss grpc pool
	faissGrpcConf := dataConf["faissGrpcAddr"].(map[string]interface{})
	faissGrpcPool, err := utils.CreateGrpcConn(faissGrpcConf)
	if err != nil {
		return err
	}

	indexInfo := dataConf["indexInfo"].(map[string]interface{})
	for indexName, tmpIndexConf := range indexInfo { //only 1 index
		tmpIndexConfMap := tmpIndexConf.(map[string]interface{})
		recallNum := int32(tmpIndexConfMap["recall_num"].(float64))
		signature := tmpIndexConfMap["index_name"].(string)
		indexInfoStruct := &faiss_index.RecallRequest{
			IndexName: signature,
			RecallNum: recallNum,
		}

		f.setIndexName(indexName)
		f.setFaissGrpcPool(faissGrpcPool)
		f.setFaissIndexs(indexInfoStruct)
	}

	return nil
}
