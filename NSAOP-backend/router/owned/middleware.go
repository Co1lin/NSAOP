package owned

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils/constant"
)

func GetUserOwnedDataByID(c *gin.Context, id interface{}, associaton string, modelArray interface{}) error {
	var user model.User
	if err := model.DB.First(&user, c.MustGet("userId").(uint)).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, constant.MsgTokenCannotFindUser)
		return err
	}
	if model.DB.Model(&user).Where(constant.IdQuery, id).Association(associaton).Count() != 1 {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgIdNotAvailable)
		return fmt.Errorf(constant.MsgIdNotAvailable)
	}
	model.DB.Model(&user).Where(constant.IdQuery, id).Association(associaton).Find(modelArray)
	return nil
}
