package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nsaop/model"
	"nsaop/router/jwt"
	"nsaop/router/resp"
)

// @Tags		User
// @Summary		refresh
// @Description	given an vaild refresh_token in cookie to gain new access_token by body
// @Product		json
// @Success		200 {object} resp.Response{data=UserToken} "ok"
// @Failure		401 {object} resp.Response "token not found<br>token has banned<br>[otherMsg]"
// @Router		/user/refresh [POST]
func Refresh(c *gin.Context) {
	if refreshToken, err := c.Cookie("refresh_token"); err == nil {
		j := jwt.NewJWT("refresh")
		claims, e := j.ParseToken(refreshToken)
		if e != nil {
			resp.ERROR(c, http.StatusUnauthorized, e.Error())
		} else {
			var r model.Refresh
			if err := model.DB.Where("token = ?", refreshToken).First(&r).Error; err != nil {
				resp.ERROR(c, http.StatusUnauthorized, "token has banned")
			} else {
				newToken := jwt.GenerateTokenForUser(c, claims.UserId, claims.UserRole, "access")
				resp.OK(c, http.StatusOK, gin.H{
					"token": newToken,
				})
			}
		}
	} else {
		resp.ERROR(c, http.StatusUnauthorized, "token not found")
	}
}
