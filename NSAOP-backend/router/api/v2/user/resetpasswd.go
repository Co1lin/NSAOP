package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"nsaop/config"
	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils"
	"nsaop/utils/constant"
	"nsaop/utils/email"
	"time"
)

type PasswdResetRequest struct {
	Username	string 	`json:"username" form:"username"`
	ReCAPTCHA	string 	`json:"g-recaptcha-response" form:"g-recaptcha-response"`
}

type ResetPasswdObj struct {
	Username	string 	`json:"username"`
	Token		string	`json:"token"`
	Password	string	`json:"password"`
}

// @Tags		User
// @Summary		Request Password Reset
// @Description	Need a username and reCAPTCHA response; the host will send an email to the user
// @Accept		json
// @Product		json
// @Param		g-recaptcha-response query string true "response token got from google reCAPTCHA (expected action: request_passwd_reset)"
// @Param		username query string true "username"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "reCAPTCHA failed<br>format error<br>field empty<br>permission denied<br>invalid info"
// @Failure		500 {object} resp.Response ""
// @Router		/user/resetpasswd [GET]
func ResetPasswdRequest(c *gin.Context) {
	var req PasswdResetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if reCAPTCHASucceed, err := utils.ReCAPTCHA(req.ReCAPTCHA, "request_passwd_reset", c.ClientIP()); !reCAPTCHASucceed {
		log.Println(err)
		resp.ERROR(c, http.StatusForbidden, "reCAPTCHA failed")
		return
	}
	if !model.ValidUsernameLen(req.Username) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	var user model.User
	if err := model.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		resp.ERROR(c, http.StatusBadRequest, "user not found")
		return
	}
	if time.Now().Sub(user.LastOper).Minutes() < 1 {
		resp.ERROR(c, http.StatusBadRequest, "too frequent request")
		return
	}
	token := utils.GetRandomStringWithLetters(32)
	subject := "NSAOP Reset Password 重置密码"
	url := "https://" + config.FrontendDomain + "/reset?username=" + user.Username + "&token=" + token
	body := fmt.Sprintf("%s:<br><br>您可以访问 <a href=\"%s\">%s</a> 重置密码。请在10分钟内完成操作。<br><br>NSAOP", user.Username, url, url)
	if err := email.SendEmailToUser(user, subject, body);
	err != nil {
		resp.ERROR(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		user.ResetToken = token
		user.LastOper = time.Now()
		model.DB.Save(&user)
		resp.OK(c, http.StatusOK, constant.MsgOK)
		return
	}
}

// @Tags		User
// @Summary		Reset Password
// @Description	Reset password
// @Accept		json
// @Product		json
// @Param		body body ResetPasswdObj true "username + new passowrd + token"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>field empty<br>invalid info"
// @Failure		500 {object} resp.Response "format error<br>field empty<br>invalid info"
// @Router		/user/resetpasswd [POST]
func ResetPassword(c *gin.Context) {
	var resetPasswd ResetPasswdObj
	if err := c.ShouldBindJSON(&resetPasswd); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if !model.ValidUsernameLen(resetPasswd.Username) ||
		!model.ValidPasswordLen(resetPasswd.Password) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	var user model.User
	if err := model.DB.Where("username = ?", resetPasswd.Username).First(&user).Error; err != nil {
		resp.ERROR(c, http.StatusBadRequest, "user not found")
		return
	}
	if resetPasswd.Token != user.ResetToken {
		resp.ERROR(c, http.StatusBadRequest, "invalid token")
		return
	}
	if time.Now().Sub(user.LastOper).Minutes() > 10 {
		resp.ERROR(c, http.StatusBadRequest, "token timeout")
		return
	}
	user.ResetToken = ""
	user.LastOper = time.Unix(0, 0)
	hash, _ := bcrypt.GenerateFromPassword([]byte(resetPasswd.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	model.DB.Save(&user)
	resp.OK(c, http.StatusOK, constant.MsgOK)
	return
}