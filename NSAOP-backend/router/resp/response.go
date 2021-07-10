package resp

import (
	"github.com/gin-gonic/gin"

	"nsaop/utils/constant"
)

type Response struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func JSON(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(code, Response{
		Data: data,
		Msg:  msg,
	})
}

func OK(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Data: data,
		Msg:  constant.MsgOK,
	})
}

func ERROR(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Data: gin.H{},
		Msg:  msg,
	})
}
