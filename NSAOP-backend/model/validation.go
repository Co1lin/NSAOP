package model

import (
	"fmt"
	"regexp"
	"time"

	"nsaop/utils/constant"
)

func ValidPayType(paytype string) bool {
	return paytype == constant.PaytypeYear || paytype == constant.PaytypeMonth
}

func ValidStatus(status string) bool {
	return status == constant.StatusWaiting || status == constant.StatusPass || status == constant.StatusOn || status == constant.StatusPause || status == constant.StatusRetrieve || status == constant.StatusCanceled || status == constant.StatusSuspend
}

func ValidRequire(require int) bool {
	return 0 <= require && require <= 7
}

func Time2String(myTime time.Time) (timeStr string) {
	timeStr = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		myTime.Year(), myTime.Month(),
		myTime.Day(), myTime.Hour(),
		myTime.Minute(), myTime.Second(),
	)
	return
}

func ValidCommentLen(comment string) bool {
	return len(comment) <= 30
}

func ValidDetailLen(detail string) bool {
	return len(detail) <= 600
}

func ValidAddressLen(site string) bool {
	return len(site) <= 300
}

func ValidContactLen(phone string) bool {
	return len(phone) <= 60
}

func ValidRole(role string) bool {
	return role == constant.RoleCustomer || role == constant.RoleOperator || role == constant.RoleEngineer || role == constant.RoleAdmin
}

func ValidUsernameLen(username string) bool {
	return 5 <= len(username) && len(username) <= 60
}

func ValidPasswordLen(password string) bool {
	return len(password) <= 32
}

func ValidCompanyLen(company string) bool {
	return 0 < len(company) && len(company) <= 90
}

func ValidPhoneLen(phone string) bool {
	return len(phone) <= 20
}

func ValidEmailLen(email string) bool {
	return len(email) <= 40
}

var validEmailEpr = regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)

func ValidEmailExpr(email string) bool {
	return validEmailEpr.MatchString(email)
}

func ValidEmail(email string) bool {
	return ValidEmailLen(email) && ValidEmailExpr(email)
}
