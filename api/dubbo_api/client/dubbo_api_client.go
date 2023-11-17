package client

import (
	"context"
	dubbo_api "infer-microservices/api/dubbo_api/server"
	"infer-microservices/pkg/logs"
	"infer-microservices/pkg/services/io"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

type DubboApiClient struct {
	dubboConfigFile string
}

func (d *DubboApiClient) setDubboConfigFile(dubboConfigFile string) {
	d.dubboConfigFile = dubboConfigFile
}

func (c *DubboApiClient) dubboServiceApiInfer(req io.RecRequest) (*io.RecResponse, error) {
	// export DUBBO_GO_CONFIG_PATH=dubbogo.yml or load it in code.
	if err := config.Load(config.WithPath(c.dubboConfigFile)); err != nil {
		panic(err)
	}

	// request dubbo infer service. use to test serivce.
	rsp, err := dubbo_api.DubboServiceApiClient.RecommenderInfer(context.TODO(), &req)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return rsp, nil
}
