package v1

import "github.com/gin-gonic/gin"

type Router struct {
	*gin.Engine
}

func Init(engine *gin.Engine) {
	r := &Router{Engine: engine}
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", r.Test)
		v1.POST("/tx/pending", r.PendingTransaction)
	}
}
