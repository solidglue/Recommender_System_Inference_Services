package nacos_config_listener

import (
	"encoding/json"
	"infer-microservices/utils"
)

type nacosContent struct {
	author  string
	update  string
	version string
	config  Config_
}

type Config_ struct {
	businessdomain string                 //share redis by domain.
	redisConfNacos map[string]interface{} //features redis conf.
	modelConfNacos map[string]interface{} //model trainning and model infer conf.
	indexConfNacos map[string]interface{} //faiss index conf.
}

// parse service config file, which contains index info、redis info and model info etc.
func (s *nacosContent) InputServiceConfigParse(content string) (string, string, string, string) {
	json.Unmarshal([]byte(string(content)), s)
	redisConfStr := utils.ConvertStructToJson(s.config.redisConfNacos)
	modelConfStr := utils.ConvertStructToJson(s.config.modelConfNacos)
	indexConfStr := utils.ConvertStructToJson(s.config.indexConfNacos)
	business := s.config.businessdomain

	return business, redisConfStr, modelConfStr, indexConfStr
}
