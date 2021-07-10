package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils/constant"
)

type UserModify struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
	Company     string `json:"company"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}

type UserDetail struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	Company  string `json:"company"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func Detail(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		GetDetail(c)
	case http.MethodPut:
		ModifyDetail(c)
	default:
		resp.ERROR(c, http.StatusMethodNotAllowed, constant.MsgMethodNotAllowed)
	}
}

// @Tags		User
// @Summary		detail
// @Description	with jwt auth, get user detail information
// @Accept		json
// @Product		json
// @Param 		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response{data=UserDetail}
// @Failure		401 {object} resp.Response "authHeader not found<br>authHeader format error<br>Token is expired<br>Couldn't handle this token<br>[otherMsg]"
// @Router		/user/detail [GET]
func GetDetail(c *gin.Context) {
	var u model.User
	model.DB.First(&u, c.MustGet("userId").(uint))
	resp.OK(c, http.StatusOK, UserDetail{
		Role:     u.Role,
		Username: u.Username,
		Company:  u.Company,
		Phone:    u.Phone,
		Email:    u.Email,
	})
}

// @Tags		User
// @Summary		modify
// @Description	with jwt auth, modify user information
// @Description	*can be null:* {
// @Description   old_passowrd, // needed when password is not empty
// @Description   password,
// @Description   company,
// @Description   phone,
// @Description   email,
// @Description }
// @Accept		json
// @Product		json
// @Param Authorization header string true "Bearer [access_token]"
// @Param		body body UserModify true "modified information."
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>invalid info"
// @Router		/user/detail [PUT]
func ModifyDetail(c *gin.Context) {
	var user UserModify
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	// check length (drop toooo long info)
	if (user.OldPassword != "" && !model.ValidPasswordLen(user.OldPassword)) ||
		(user.Password != "" && !model.ValidPasswordLen(user.Password)) ||
		(user.Company != "" && !model.ValidCompanyLen(user.Company)) ||
		(user.Email != "" && !model.ValidEmail(user.Email)) ||
		(user.Phone != "" && !model.ValidPhoneLen(user.Phone)) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	var u model.User
	if err := model.DB.First(&u, c.MustGet("userId").(uint)).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, constant.MsgTokenCannotFindUser)
		return
	}
	if user.Password != "" {
		if err := bcrypt.CompareHashAndPassword(
			[]byte(u.Password), []byte(user.OldPassword)); err != nil {
			resp.ERROR(c, http.StatusUnauthorized, constant.MsgWrongPassword)
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		u.Password = string(hash)
	}
	if user.Company != "" {
		u.Company = user.Company
	}
	if user.Phone != "" {
		u.Phone = user.Phone
	}
	if user.Email != "" {
		u.Email = user.Email
	}
	model.DB.Save(&u)
	resp.OK(c, http.StatusOK, gin.H{})
}
