package infer_pipeline

import (
	config_loader "infer-microservices/pkg/config_loader"
	"infer-microservices/pkg/infer_samples/feature"

	"net/http"
	"strings"
)

//https://developer.aliyun.com/article/1065404

//把样本、召回、粗排、精排、重排等串成一个工作流

//各个模型需要实现统一的接口，infer。服务通过pipline调用模型

type Pipeline struct {
	steps []inferAlgMap
}

type inferAlgMap struct {
	algName string
	algFunc InferPipelineInterface //sample class ,or model class
}

// TODO: 将steps字符串转化为steps
func (p Pipeline) SetSteps(steps []string) {
	//p.steps = steps
}

// TODO:策略模式
// steps :[("sample",inferSampleDirector.RecallSampleDirector),("recall",RecallSampleDirector),("pre_ranking",RecallSampleDirector),("re_rank",RecallSampleDirector)],
func (p Pipeline) Predict(serviceConfig *config_loader.ServiceConfig, requestId string, userId string, r *http.Request, lightInfer bool) (map[string][]map[string]interface{}, error) {
	var err error
	recallSample := feature.ExampleFeatures{}
	preRankingSample := feature.ExampleFeatures{}
	rankingSample := feature.ExampleFeatures{}

	recallResponse := make(map[string][]map[string]interface{}, 0)
	prerankingResponse := make(map[string][]map[string]interface{}, 0)
	rankingResponse := make(map[string][]map[string]interface{}, 0)
	reRankResponse := make(map[string][]map[string]interface{}, 0)

	recallItemIdList := make([]string, 0)
	preRankingItemIdList := make([]string, 0)
	rankingItemIdList := make([]string, 0)

	for index, step := range p.steps {
		//TODO:所有的推理模型和样本，需要实现一个infer接口

		//1.RECALL
		//recall sample
		if strings.Contains(step.algName, "recall_sample") {
			emptIdList := make([]string, 0)
			recallSample, err = p.inferSample(serviceConfig, step, requestId, userId, r, emptIdList)
			if err != nil {
				return reRankResponse, err
			}
		}
		//recall infer
		recallResponse, err = p.recall(serviceConfig, step, requestId, r, recallSample, lightInfer)
		//return the last model in pipeline.
		if index == len(p.steps)-1 {
			return recallResponse, nil
		}

		//2.PRE-RANKING
		//recall items
		for _, itemid := range recallResponse["data"] {
			recallItemIdList = append(recallItemIdList, itemid["itemid"].(string))
		}
		//prerank samples
		if strings.Contains(step.algName, "preranking_sample") {
			preRankingSample, err = p.inferSample(serviceConfig, step, requestId, userId, r, recallItemIdList)
			if err != nil {
				return reRankResponse, err
			}
		}
		//prerank infer
		prerankingResponse, err = p.rank(serviceConfig, step, requestId, r, preRankingSample, lightInfer)
		//return the last model in pipeline.
		if index == len(p.steps)-1 {
			return prerankingResponse, nil
		}

		//3.RANKING
		//preranking items
		for _, itemid := range prerankingResponse["data"] {
			preRankingItemIdList = append(preRankingItemIdList, itemid["itemid"].(string))
		}
		//prerank samples
		if strings.Contains(step.algName, "ranking_sample") && !strings.Contains(step.algName, "preranking_sample") {
			rankingSample, err = p.inferSample(serviceConfig, step, requestId, userId, r, recallItemIdList)
			if err != nil {
				return reRankResponse, err
			}
		}
		//prerank infer
		rankingResponse, err = p.rank(serviceConfig, step, requestId, r, rankingSample, lightInfer) //return all the preranking items
		//return the last model in pipeline.
		if index == len(p.steps)-1 {
			return rankingResponse, nil
		}

		//4.RE-RANK by rules
		//ranking items
		for _, itemid := range rankingResponse["data"] {
			rankingItemIdList = append(rankingItemIdList, itemid["itemid"].(string))
		}
		reRankResponse, err = p.re_rank(step, requestId, rankingItemIdList)
		if err != nil {
			return reRankResponse, err
		}
		//return the last model in pipeline.
		if index == len(p.steps)-1 {
			return rankingResponse, nil
		}

	}
	return reRankResponse, nil
}

// SAMPLE
func (p Pipeline) inferSample(serviceConfig *config_loader.ServiceConfig, step inferAlgMap, requestId string, userId string, r *http.Request, itemIdList []string) (feature.ExampleFeatures, error) {
	var err error
	inferSample := feature.ExampleFeatures{}
	modelsConf := serviceConfig.GetModelsConfig()
	modleConf := modelsConf[step.algName]
	offlineFeature := true
	onlineFeature := true
	featureList := modleConf.GetFeatureList()
	if modleConf.GetUserRedisKeyPreOffline() == "" {
		offlineFeature = false
	}
	if modleConf.GetUserRedisKeyPreRealtime() == "" {
		offlineFeature = false
	}

	if strings.Contains(step.algName, "recall") {
		inferSample, err = step.algFunc.RecallSampleDirector(userId, offlineFeature, onlineFeature, featureList)
		if err != nil {
			return inferSample, err
		}
	} else if strings.Contains(step.algName, "ranking") {
		inferSample, err = step.algFunc.RankingSampleDirector(userId, offlineFeature, onlineFeature, itemIdList, featureList)
		if err != nil {
			return inferSample, err
		}
	} else {
		// support more samples.
	}

	return inferSample, nil
}

// RECALL  TFSERVING INFER
func (p Pipeline) recall(serviceConfig *config_loader.ServiceConfig, step inferAlgMap, requestId string, r *http.Request, recallSample feature.ExampleFeatures, lightInfer bool) (map[string][]map[string]interface{}, error) {
	var err error
	recallResponse := make(map[string][]map[string]interface{}, 0)
	modelsConf := serviceConfig.GetModelsConfig()
	pipelineConf := serviceConfig.GetPipelineConfig()

	recallNum := 100
	if lightInfer {
		recallNum = int(pipelineConf.GetRecallNum())
	} else {
		recallNum = int(pipelineConf.GetRecallNumLight())
	}

	if strings.Contains(step.algName, "recall_skywalking") {
		recallResponse, err = step.algFunc.ModelInferSkywalking(modelsConf[step.algName], requestId, r, recallSample, recallNum)
		if err != nil {
			return recallResponse, err
		}

	} else if strings.Contains(step.algName, "recall_noskywalking") {
		recallResponse, err = step.algFunc.ModelInferNoSkywalking(modelsConf[step.algName], requestId, recallSample, recallNum)
		if err != nil {
			return recallResponse, err
		}
	} else {
		// support more models.
	}

	return recallResponse, nil
}

// RANKING，RANKING  TFSERVING INFER
func (p Pipeline) rank(serviceConfig *config_loader.ServiceConfig, step inferAlgMap, requestId string, r *http.Request, recallSample feature.ExampleFeatures, lightInfer bool) (map[string][]map[string]interface{}, error) {
	var err error
	recallResponse := make(map[string][]map[string]interface{}, 0)
	modelsConf := serviceConfig.GetModelsConfig()
	pipelineConf := serviceConfig.GetPipelineConfig()

	rankNum := 100
	if lightInfer {
		rankNum = int(pipelineConf.GetPreRankingNumLight())
	} else {
		rankNum = int(pipelineConf.GetPreRankingNumLight())
	}

	if strings.Contains(step.algName, "ranking_skywalking") {
		recallResponse, err = step.algFunc.ModelInferSkywalking(modelsConf[step.algName], requestId, r, recallSample, rankNum)
		if err != nil {
			return recallResponse, err
		}

	} else if strings.Contains(step.algName, "ranking_noskywalking") {
		recallResponse, err = step.algFunc.ModelInferNoSkywalking(modelsConf[step.algName], requestId, recallSample, rankNum)
		if err != nil {
			return recallResponse, err
		}
	} else {
		// support more models.
	}

	return recallResponse, nil
}

func (p Pipeline) re_rank(step inferAlgMap, requestId string, rankingItemIdList []string) (map[string][]map[string]interface{}, error) {
	var err error
	reRankResponse := make(map[string][]map[string]interface{}, 0)

	return reRankResponse, err
}
