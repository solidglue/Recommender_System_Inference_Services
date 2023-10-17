package apis

import "infer-microservices/cores/service_config_loader"

var ServiceConfigs = make(map[string]*service_config_loader.ServiceConfig, 0) //one server/dataid,one service conn
var NacosListedMap = make(map[string]bool, 0)
