package routers

import (
	"github.com/gin-gonic/gin"
	"vbh/btc-plugins/api/routers/v1"
)

func Init() *gin.Engine {
	r := gin.Default()
	v1.Init(r)
	return r
}
