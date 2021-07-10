package user_test

import (
	"net/http"
	"testing"

	"nsaop/router/api/v2/user"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestCheck(t *testing.T) {
	url := "/user/check/username"

	goodCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: test.CommonUsers[constant.RoleCustomer].Username,
			},
			Code: http.StatusConflict,
			Msg:  "username exists",
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: test.CommonUsers[constant.RoleOperator].Username,
			},
			Code: http.StatusConflict,
			Msg:  "username exists",
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: test.CommonUsers[constant.RoleEngineer].Username,
			},
			Code: http.StatusConflict,
			Msg:  "username exists",
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: test.CommonUsers[constant.RoleAdmin].Username,
			},
			Code: http.StatusConflict,
			Msg:  "username exists",
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: "unknown",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
		},
	}

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})

	badCases := []test.TestCase{
		{
			Desc:   "Format error",
			Url:    url,
			Method: http.MethodPost,
			Param: struct {
				Username int `json:"username"`
			}{
				Username: 1,
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFormatError,
		},
		{
			Desc:   "Username too short",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: "NSAO",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   "Username too long",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserCheck{
				Username: "这是一个二十字的名字这是一个二十字的名字超了",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
	}

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, badCases, false)
	})
}

func TestCheckPassword(t *testing.T) {
	url := "/user/check/password"
	token, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	goodCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			Param: user.PasswordCheck{
				Password: "TestPasswd@nsaop2021",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
		},
	}

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})

	badCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			Param: user.PasswordCheck{
				Password: "WrongPasswd@nsaop2021",
			},
			Code: http.StatusUnauthorized,
			Msg:  constant.MsgWrongPassword,
		},
		{
			Desc:   "Format Error",
			Url:    url,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			Param: struct {
				Password int `json:"password"`
			}{
				Password: 123,
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFormatError,
		},
		{
			Desc:   "Wrong Password Len",
			Url:    url,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			Param: user.PasswordCheck{
				Password: "TestPasswd@nsaop2021TestPasswd@nsaop2021TestPasswd@nsaop2021TestPasswd@nsaop2021",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
	}

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, badCases, false)
	})
}
