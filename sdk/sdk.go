package sdk

import (
	"github.com/gin-gonic/gin"

	"gucooing/lolo/gdconf"
	"gucooing/lolo/pkg/flyrsa"
	"gucooing/lolo/pkg/log"
)

type Server struct {
	router *gin.Engine
	rsa    *flyrsa.PrivateKey
}

func New(router *gin.Engine) *Server {
	s := &Server{
		router: router,
	}
	s.Router()
	priv, err := flyrsa.NewPrivateKey(gdconf.GetPrivatePem())
	if err != nil {
		log.App.Warnf("初始化FlyRsa失败:%s", err.Error())
	}
	s.rsa = priv
	return s
}

func (s *Server) Close() {}
