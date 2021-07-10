package service

import (
	"net/http"
	"nsaop/model"
	"nsaop/router/owned"
	"nsaop/router/resp"
	"nsaop/utils/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func IdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := uuid.Parse(c.Param("id")); err != nil {
			resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
			c.Abort()
		} else {
			var services []model.Service
			if err := owned.GetUserOwnedDataByID(c, id, constant.TableService, &services); err != nil {
				c.Abort()
			} else {
				c.Set(constant.TableService, services[0])
				c.Next()
			}
		}
	}
}
