package nacos_config_listener

import (
	"encoding/json"
	"infer-microservices/utils"
	"infer-microservices/utils/logs"

	validator "github.com/go-playground/validator/v10"
)

type nacosContent struct {
	author  string  `validate:"required"`
	update  string  `validate:"required"`
	version string  `validate:"required"`
	config  Config_ `validate:"required"`
}

type Config_ struct {
	redisConfNacos map[string]interface{} `validate:"required"` //features redis conf.
	modelConfNacos map[string]interface{} `validate:"required"` //model trainning and model infer conf.
	indexConfNacos map[string]interface{} //faiss index conf.
}

// parse service config file, which contains index info、redis info and model info etc.
func (s *nacosContent) InputServiceConfigParse(content string) (string, string, string) {

	json.Unmarshal([]byte(string(content)), s)
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		logs.Error(err)
		return "", "", ""
	}

	redisConfStr := utils.ConvertStructToJson(s.config.redisConfNacos)
	modelConfStr := utils.ConvertStructToJson(s.config.modelConfNacos)
	indexConfStr := utils.ConvertStructToJson(s.config.indexConfNacos)

	return redisConfStr, modelConfStr, indexConfStr
}
