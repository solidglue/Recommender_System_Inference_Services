package nacos

import (
	"encoding/json"
	"infer-microservices/internal/logs"
	"infer-microservices/internal/utils"

	validator "github.com/go-playground/validator/v10"
)

type NacosContent struct {
	// author  string  `validate:"required"`
	// update  string  `validate:"required"`
	// version string  `validate:"required"`
	Config Config_ `validate:"required"`
}

type Config_ struct {
	redisConfNacos map[string]interface{} `validate:"required"` //features redis conf.
	modelConfNacos map[string]interface{} `validate:"required"` //model trainning and model infer conf.
	indexConfNacos map[string]interface{} //faiss index conf.
}

// parse service config file, which contains index info、redis info and model info etc.
func (s *NacosContent) InputServiceConfigParse(content string) (string, string, string) {
	tmpNacos := &NacosContent{}
	json.Unmarshal([]byte(string(content)), tmpNacos)
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		logs.Error(err)
		return "", "", ""
	}

	redisConfStr := utils.ConvertStructToJson(s.Config.redisConfNacos)
	modelConfStr := utils.ConvertStructToJson(s.Config.modelConfNacos)
	indexConfStr := utils.ConvertStructToJson(s.Config.indexConfNacos)

	return redisConfStr, modelConfStr, indexConfStr
}
