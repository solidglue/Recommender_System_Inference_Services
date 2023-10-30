package nacos_config_listener

import (
	"encoding/json"
	"infer-microservices/utils"
)

var nacosContentInstance *NacosContent

type NacosContent struct {
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

// func init() {
// 	//serviceConfFile = *flags.Service_start_file
// 	nacosContentInstance = new(NacosContent)

// }

// func getNacosContentInstance() *NacosContent {
// 	return nacosContentInstance
// }

// parse service config file, which contains index info、redis info and model info etc.
func (s *NacosContent) InputServiceConfigParse(content string) (string, string, string, string) {
	json.Unmarshal([]byte(string(content)), s)
	redisConfStr := utils.ConvertStructToJson(s.config.redisConfNacos)
	modelConfStr := utils.ConvertStructToJson(s.config.modelConfNacos)
	indexConfStr := utils.ConvertStructToJson(s.config.indexConfNacos)
	business := s.config.businessdomain

	return business, redisConfStr, modelConfStr, indexConfStr
}
