package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nsaop/router/resp"
	"nsaop/utils/constant"
	"nsaop/utils/email"
)

// @Tags		User
// @Summary		Send email to the requested user
// @Description	with jwt auth, send email to the user
// @Accept		json
// @Product		json
// @Param 		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "permission denied"
// @Router		/user/email [GET]
func SendEmail(c *gin.Context) {
	if c.MustGet("userRole") != constant.RoleOperator {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}
	resp.OK(c, http.StatusOK, gin.H{})
	email.SendBillAll()
}
