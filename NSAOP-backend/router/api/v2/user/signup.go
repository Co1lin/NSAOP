package user

import (
	"log"
	"net/http"
	"nsaop/config"
	"nsaop/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils/constant"
)

type UserSignup struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	Password string `json:"password"`
	Company  string `json:"company"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Code	 string	`json:"code"`
	ReCAPTCHA	string `json:"g_recaptcha_response"`
}

// @Tags		User
// @Summary		signup
// @Description	signup with user information
// @Description	**must fill:**
// @Description <pre>{
// @Description   role, // options: customer, operator, engineer, admin
// @Description   username, // (5 <= length <= 20)
// @Description   password, // (length == 32)
// @Description   company,	// (0 < length <= 30)
// @Description }</pre>
// @Description	**can be null:**
// @Description	<pre>{
// @Description   phone, // (length <= 20)
// @Description   email, // (length <= 40)
// @Description   code, // (length <= 40)
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param		g-recaptcha-response query string true "response token got from google reCAPTCHA (expected action: register)"
// @Param		body body UserSignup true "user information."
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "reCAPTCHA failed<br>format error<br>field empty<br>invalid user role<br>invalid info (<username, password, company, phone, email>)<br>username already exist"
// @Router		/user/signup [POST]
func Signup(c *gin.Context) {
	var user UserSignup
	if err := c.ShouldBindJSON(&user); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if reCAPTCHASucceed, err := utils.ReCAPTCHA(user.ReCAPTCHA, "register", c.ClientIP()); !reCAPTCHASucceed {
		log.Println(err)
		resp.ERROR(c, http.StatusForbidden, "reCAPTCHA failed")
		return
	}
	if user.Role == "" ||
		user.Username == "" ||
		user.Password == "" ||
		user.Company == "" {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFieldEmpty)
		return
	}
	if !model.ValidRole(user.Role) ||
		!model.ValidUsernameLen(user.Username) ||
		!model.ValidPasswordLen(user.Password) ||
		!model.ValidCompanyLen(user.Company) ||
		!model.ValidPhoneLen(user.Phone) ||
		(user.Email != "" && !model.ValidEmail(user.Email)) ||
		(user.Role != constant.RoleCustomer && user.Code != config.AdminRegisterCode) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	u := model.User{
		Role:     user.Role,
		Username: user.Username,
		Password: user.Password,
		Company:  user.Company,
		Phone:    user.Phone,
		Email:    user.Email,
	}
	if err := model.DB.Where("username = ?", u.Username).First(&u).Error; err == nil {
		resp.ERROR(c, http.StatusBadRequest, "username already exist")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
	u.LastOper = time.Unix(0, 0)

	if err := model.DB.Create(&u).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, "unknown error occur when inserting into database")
	} else {
		resp.OK(c, http.StatusOK, gin.H{})
	}
}
