package server

import (
	"infer-microservices/apis/io"

	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"
)

type DubboServer struct {
	nacosIp        string
	nacosPort      uint
	lowerRankNum   int
	lowerRecallNum int
	dubboConfFile  string
}

func init() {
	//regisger dubbo service.
	hessian.RegisterPOJO(&io.RecRequest{})
	hessian.RegisterPOJO(&io.RecResponse{})
}

//set func

func (s *DubboServer) SetNacosIp(nacosIp string) {
	s.nacosIp = nacosIp
}

func (s *DubboServer) SetNacosPort(nacosPort uint) {
	s.nacosPort = nacosPort
}

func (s *DubboServer) SetLowerRankNum(lowerRankNum int) {
	s.lowerRankNum = lowerRankNum
}

func (s *DubboServer) SetLowerRecallNum(lowerRecallNum int) {
	s.lowerRecallNum = lowerRecallNum
}

func (s *DubboServer) SetDubboConfFile(dubboConfFile string) {
	s.dubboConfFile = dubboConfFile
}

// @implement start infertace
func (s *DubboServer) ServerStart() {
	config.SetProviderService(s)
	if err := config.Load(config.WithPath(s.dubboConfFile)); err != nil {
		panic(err)
	}
}
