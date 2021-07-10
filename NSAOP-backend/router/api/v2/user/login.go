package user

import (
	"log"
	"net/http"
	"nsaop/utils"
	"nsaop/utils/constant"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"nsaop/config"
	"nsaop/model"
	"nsaop/router/jwt"
	"nsaop/router/resp"
)

type UserLogin struct {
	Username 	string `json:"username"`
	Password 	string `json:"password"`
	ReCAPTCHA	string `json:"g_recaptcha_response"`
}

type UserToken struct {
	Token string `json:"token"`
}

// @Tags		User
// @Summary		login
// @Description	login and return access_token by body and refresh_token by cookie
// @Accept		json
// @Product		json
// @Param		g-recaptcha-response query string true "response token got from google reCAPTCHA (expected action: login)"
// @Param		body body UserLogin true "username + password"
// @Success		200 {object} resp.Response{data=UserToken} "ok"
// @Failure		400 {object} resp.Response "format error<br>field empty<br>invalid info"
// @Failure		401 {object} resp.Response "user not found<br>wrong password"
// @Router		/user/login [POST]
func Login(c *gin.Context) {
	var user UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if reCAPTCHASucceed, err := utils.ReCAPTCHA(user.ReCAPTCHA, "login", c.ClientIP()); !reCAPTCHASucceed {
		log.Println(err)
		resp.ERROR(c, http.StatusForbidden, "reCAPTCHA failed")
		return
	}
	if user.Username == "" || user.Password == "" {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFieldEmpty)
		return
	}
	if !model.ValidUsernameLen(user.Username) || !model.ValidPasswordLen(user.Password) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	var u model.User
	if err := model.DB.Where("username = ?", user.Username).First(&u).Error; err != nil {
		resp.ERROR(c, http.StatusUnauthorized, "user not found")
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
			resp.ERROR(c, http.StatusUnauthorized, constant.MsgWrongPassword)
		} else {
			accessToken := jwt.GenerateTokenForUser(c, u.ID, u.Role, "access")
			refreshToken := jwt.GenerateTokenForUser(c, u.ID, u.Role, "refresh")

			if *config.IfCORS {
				c.SetSameSite(4)
			}

			c.SetCookie("refresh_token", refreshToken, config.Router.GetInt("jwt.refreshTime"), "", "", true, true)

			model.DB.Create(&model.Refresh{Token: refreshToken, CreateAt: time.Now()})
			resp.OK(c, http.StatusOK, UserToken{
				Token: accessToken,
			})
		}
	}
}
