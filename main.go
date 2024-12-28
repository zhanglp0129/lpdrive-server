package main

import (
	"fmt"
	"github.com/zhanglp0129/lpdrive-server/config"
	"github.com/zhanglp0129/lpdrive-server/logger"
	"github.com/zhanglp0129/lpdrive-server/router"
	"net/http"
)

func main() {
	ip, port := config.C.Server.IP, config.C.Server.Port
	address := fmt.Sprintf("%s:%d", ip, port)
	// 运行Web服务
	err := http.ListenAndServe(address, router.R)
	if err != nil {
		logger.L.WithError(err).Panicln("Web服务器启动失败")
	}
}
