package pipeline_config

import (
	"infer-microservices/internal/utils"
)

type PipelineConfig struct {
	recallNum          int32 //recall num
	preRankingNum      int32 //filter recall items, use to ranking
	recallNumLight     int32 //recall num ,visit it when service degradation
	preRankingNumLight int32 //filter recall items, use to ranking.visit it when service degradation
	pipeline           []string
	lightPipeline      []string
}

// recallNum
func (r *PipelineConfig) SetRecallNum(recallNum int32) {
	r.recallNum = recallNum
}

func (r *PipelineConfig) GetRecallNum() int32 {
	return r.recallNum
}

// preRankingNum
func (r *PipelineConfig) SetPreRankingNum(preRankingNum int32) {
	r.preRankingNum = preRankingNum
}

// recallNumLight
func (r *PipelineConfig) SetRecallNumLight(recallNumLight int32) {
	r.recallNumLight = recallNumLight
}

func (r *PipelineConfig) GetRecallNumLight() int32 {
	return r.recallNumLight
}

// preRankingNumLight
func (r *PipelineConfig) SetPreRankingNumLight(preRankingNumLight int32) {
	r.preRankingNumLight = preRankingNumLight
}

func (r *PipelineConfig) GetPreRankingNumLight() int32 {
	return r.preRankingNumLight
}

// pipeline
func (r *PipelineConfig) SetPipeline(pipeline []string) {
	r.pipeline = pipeline
}

func (r *PipelineConfig) GetPipeline() []string {
	return r.pipeline
}

// lightPipeline
func (r *PipelineConfig) SetLightPipeline(lightPipeline []string) {
	r.lightPipeline = lightPipeline
}

func (r *PipelineConfig) GetLightPipeline() []string {
	return r.lightPipeline
}

// @implement ConfigLoadInterface
func (f *PipelineConfig) ConfigLoad(dataId string, pipelineConfStr string) error {
	dataConf := utils.ConvertJsonToStruct(pipelineConfStr)
	recallNum := dataConf["recallNum"].(int32)
	preRankingNum := dataConf["preRankingNum"].(int32)
	recallNumLight := dataConf["recallNumLight"].(int32)
	preRankingNumLight := dataConf["preRankingNumLight"].(int32)
	pipline := dataConf["pipeline"].([]string)
	lightPipline := dataConf["light_pipeline"].([]string)

	f.SetRecallNum(recallNum)
	f.SetPreRankingNum(preRankingNum)
	f.SetRecallNumLight(recallNumLight)
	f.SetRecallNumLight(recallNumLight)
	f.SetPreRankingNumLight(preRankingNumLight)
	f.SetLightPipeline(lightPipline)
	f.SetPipeline(pipline)

	return nil
}
