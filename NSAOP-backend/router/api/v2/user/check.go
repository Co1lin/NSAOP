package user

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils/constant"
)

type UserCheck struct {
	Username string `json:"username"`
}

type PasswordCheck struct {
	Password string `json:"password"`
}

// @Tags		User
// @Summary		check
// @Description	check whether a given username exists
// @Accept		json
// @Product		json
// @Param		body body UserCheck true "username"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>invalid info"
// @Failure		409 {object} resp.Response "username exists"
// @Router		/user/check/username [POST]
func CheckUsername(c *gin.Context) {
	var user UserCheck
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if !model.ValidUsernameLen(user.Username) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	var u model.User
	if err := model.DB.Where("username = ?", user.Username).First(&u).Error; err != nil {
		resp.OK(c, http.StatusOK, gin.H{})
	} else {
		resp.ERROR(c, http.StatusConflict, "username exists")
	}
}

// @Tags		User
// @Summary		check password
// @Description	check whether the password is correct
// @Accept		json
// @Product		json
// @Param Authorization header string true "Bearer [access_token]"
// @Param		body body PasswordCheck true "Password"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>invalid info"
// @Failure		401 {object} resp.Response "wrong password"
// @Router		/user/check/password [POST]
func CheckPassword(c *gin.Context) {
	var passwd PasswordCheck
	if err := c.ShouldBindJSON(&passwd); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if !model.ValidPasswordLen(passwd.Password) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	// get the token and bind to user
	var u model.User
	if err := model.DB.First(&u, c.MustGet("userId").(uint)).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, constant.MsgTokenCannotFindUser)
		return
	}
	// compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password), []byte(passwd.Password)); err != nil {
		resp.ERROR(c, http.StatusUnauthorized, constant.MsgWrongPassword)
	} else {
		resp.OK(c, http.StatusOK, gin.H{})
	}
}
