package apis

import "infer-microservices/cores/service_config"

var ServiceConfigs = make(map[string]*service_config.ServiceConfig, 0) //one server/dataid,one service conn
var NacosListedMap = make(map[string]bool, 0)
