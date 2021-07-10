package user_test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nsaop/model"
	"testing"
	"time"

	"nsaop/router/api/v2/user"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

var url = "/user/resetpasswd"

func TestResetPasswdRequest(t *testing.T) {
	defer test.SetupDatabaseForTest([]string{constant.TableUser, constant.TableLocation})()
	t.Run("", func(t *testing.T) {

		testUser := test.CommonUsers[constant.RoleCustomer]
		testUser.LastOper = time.Now()
		model.DB.Save(&testUser)

		badCases := []test.TestCase{
			/*{
				Desc:   "Wrong Format",
				Url:    url + "?username=1&username=2",
				Method: http.MethodGet,
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},*/
			{
				Desc:   "Invalid Username Length",
				Url:    url + "?username=f***",
				Method: http.MethodGet,
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "User Not Found",
				Url:    url + "?username=ffff***",
				Method: http.MethodGet,
				Code: http.StatusBadRequest,
				Msg:  "user not found",
			},
			{
				Desc:   "Too Frequent Request",
				Url:    url + "?username=" + testUser.Username,
				Method: http.MethodGet,
				Code: http.StatusBadRequest,
				Msg:  "too frequent request",
			},
		}
		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})
	})
}

func TestResetPassword(t *testing.T) {
	defer test.SetupDatabaseForTest([]string{constant.TableUser, constant.TableLocation})()
	t.Run("", func(t *testing.T) {

		testUser := test.CommonUsers[constant.RoleCustomer]
		testUser.LastOper = time.Unix(0, 0)
		model.DB.Save(&testUser)

		badCases := []test.TestCase{
			{
				Desc:   "Wrong Format",
				Url:    url,
				Param: gin.H{
					"username": 1,
				},
				Method: http.MethodPost,
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "Invalid Username Length",
				Url:    url,
				Method: http.MethodPost,
				Param: user.ResetPasswdObj{
					Username: "f***",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Invalid Password Length",
				Url:    url,
				Method: http.MethodPost,
				Param: user.ResetPasswdObj{
					Username: "f***",
					Password: "ffffffffffffffffffffffffffffffffffffffff",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "User Not Found",
				Url:    url,
				Method: http.MethodPost,
				Param: user.ResetPasswdObj{
					Username: "fff***",
				},
				Code: http.StatusBadRequest,
				Msg:  "user not found",
			},
			{
				Desc:   "Wrong Token",
				Url:    url,
				Method: http.MethodPost,
				Param: user.ResetPasswdObj{
					Username: testUser.Username,
					Password: "ffffff",
					Token: "ffffff",
				},
				Code: http.StatusBadRequest,
				Msg:  "invalid token",
			},
			{
				Desc:   "Timeout",
				Url:    url,
				Method: http.MethodPost,
				Param: user.ResetPasswdObj{
					Username: testUser.Username,
					Password: "ffffff",
					Token: testUser.ResetToken,
				},
				Code: http.StatusBadRequest,
				Msg:  "token timeout",
			},
		}
		t.Run("Fail", func(t *testing.T) {
			test.PerformTest(t, badCases, false)
		})

		testUser.LastOper = time.Now()
		model.DB.Save(&testUser)

		goodCases := []test.TestCase{
			{
				Desc:   "",
				Url:    url,
				Method: http.MethodPost,
				Param: user.ResetPasswdObj{
					Username: testUser.Username,
					Password: "ffffff",
					Token: testUser.ResetToken,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
		}
		t.Run("Succeed", func(t *testing.T) {
			test.PerformTest(t, goodCases, false)
		})
	})
}