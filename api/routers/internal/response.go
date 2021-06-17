package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"vbh/btc-plugins/log"
)

type Resp struct {
	Code int          `json:"code"`
	Data interface{}  `json:"data,omitempty"`
	Msg  string       `json:"msg"`
	C    *gin.Context `json:"-"`
}

func (r *Resp) Ok(data interface{}) {
	r.Data = data
	if r.Data == nil {
		r.Msg = "success"
	}
	debugLogToJson(r)
	r.C.JSON(http.StatusOK, r)
}

func (r *Resp) R(data interface{}, httpCode int) {
	r.Data = data
	debugLogToJson(r)
	r.C.JSON(httpCode, r)
}

func (r *Resp) Error(code int, msg string, err error) {
	r.Code = code
	r.Msg = msg
	debugLogToJson(r)
	log.Error(err)
	r.C.JSON(http.StatusBadRequest, r)
	r.C.Abort()
}

func (r Resp) Err(code int, err error) {
	r.Code = code
	r.Msg = err.Error()
	debugLogToJson(r)
	log.Error(err)
	r.C.JSON(http.StatusBadRequest, r)
	r.C.Abort()
}

func (r Resp) NotExist(code int, err error) {
	r.Code = code
	r.Msg = err.Error()
	debugLogToJson(r)
	log.Error(err)
	r.C.JSON(http.StatusNotFound, r)
	r.C.Abort()
}

func LogReq(c *gin.Context, req interface{}) {
	if log.IsDebug() {
		bytes, err := json.Marshal(req)
		if err != nil {
			log.Warn(err)
		}
		log.Debugf("url=%v,req=%v", c.Request.URL, string(bytes))
	}
}

func debugLogToJson(d interface{}) {
	if log.IsDebug() {
		bytes, err := json.Marshal(d)
		if err != nil {
			log.Warn(err)
		}
		log.Debug(string(bytes))
	}
}
