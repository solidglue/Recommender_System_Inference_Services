package faiss_config

import "testing"

func TestFaissConfigLoader(t *testing.T) {

	faissTestStr := `
	{
		"index_conf": {
			"faissGrpcAddr": {
				"addrs": [],
				"pool_size": 50,
				"initCap":10,
				"idleTimeoutMs": 100,
				"readTimeoutMs": 100,
				"writeTimeoutMs":100,
				"dialTimeoutS":600
			},
			"indexInfo":{
				"index-001": {
					"recallNum": 100,
					"indexName": "index-001"
				}
			}
		},
	}
	`
	faissConfs := FaissIndexConfigs{}
	faissConfs.ConfigLoad("testId", faissTestStr)

	t.Log("faissConf:", faissConfs)

	for _, faissConf := range faissConfs.faissIndexConfigs {
		if faissConf.faissIndexs.RecallNum <= 0 {
			t.Errorf("recall num must > 0")
		}

		if faissConf.faissIndexs.IndexName == "" {
			t.Errorf("index name cant be empt")
		}
	}

}
