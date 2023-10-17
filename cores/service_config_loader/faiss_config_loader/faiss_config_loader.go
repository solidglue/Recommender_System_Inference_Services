package faiss_config_loader

import (
	faiss_index "infer-microservices/common/faiss_gogofaster"
	"infer-microservices/utils"
)

var FaissIndexClientInstance *FaissIndexClient

type FaissIndexClient struct {
	indexName     string                     //index name.
	faissGrpcPool *utils.GRPCPool            //faiss  grpc pool.
	faissIndexs   *faiss_index.RecallRequest // faiss index.
}

// INFO: singleton instance
func init() {

}

// index name
func (f *FaissIndexClient) setIndexName(indexName string) {
	f.indexName = indexName
}

func (f *FaissIndexClient) GetIndexName() string {
	return f.indexName
}

// grpc pool
func (f *FaissIndexClient) setFaissGrpcPool(faissGrpcPool *utils.GRPCPool) {
	f.faissGrpcPool = faissGrpcPool
}

func (f *FaissIndexClient) GetFaissGrpcPool() *utils.GRPCPool {
	return f.faissGrpcPool
}

// FaissIndexs
func (f *FaissIndexClient) setFaissIndexs(faissIndexs *faiss_index.RecallRequest) {
	f.faissIndexs = faissIndexs
}

func (f *FaissIndexClient) GetFaissIndexs() *faiss_index.RecallRequest {
	return f.faissIndexs
}

// faiss index conf load
func (f *FaissIndexClient) ConfigLoad(domain string, dataId string, indexConfStr string) error {
	dataConf := utils.Json2Map(indexConfStr)

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
