package user

import (
	"net/http"
	"nsaop/model"
	"nsaop/router/resp"

	"github.com/gin-gonic/gin"
)

// @Tags		User
// @Summary		logout
// @Description	given an refresh_token from cookie to ban
// @Product		json
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "token not exists<br>refresh_token not found"
// @Router		/user/logout [POST]
func Logout(c *gin.Context) {
	var u model.Refresh
	if token, err := c.Cookie("refresh_token"); err == nil {
		c.SetCookie("refresh_token", token, -1, "", "", true, true)
		if e := model.DB.Where("token = ?", token).First(&u).Error; e == nil {
			model.DB.Delete(&u)
		} else {
			resp.ERROR(c, http.StatusBadRequest, "token not exists")
			return
		}
	} else {
		resp.ERROR(c, http.StatusBadRequest, "refresh_token not found")
		return
	}
	resp.OK(c, http.StatusOK, gin.H{})
}
