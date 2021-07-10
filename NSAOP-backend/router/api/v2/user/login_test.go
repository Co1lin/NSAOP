package user_test

import (
	"net/http"
	"testing"

	"nsaop/router/api/v2/user"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestLogin(t *testing.T) {
	url := "/user/login"

	goodCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: test.CommonUsers[constant.RoleCustomer].Username,
				Password: "TestPasswd@nsaop2021",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: test.CommonUsers[constant.RoleOperator].Username,
				Password: "TestPasswd@nsaop2021",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: test.CommonUsers[constant.RoleEngineer].Username,
				Password: "TestPasswd@nsaop2021",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
		},
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: test.CommonUsers[constant.RoleAdmin].Username,
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
			Desc:   "Format error",
			Url:    url,
			Method: http.MethodPost,
			Param: struct {
				Password int `json:"password"`
			}{
				Password: 123, // should be string
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFormatError,
		},
		{
			Desc:   "Missing Username & Password",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: "",
				Password: "",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFieldEmpty,
		},
		{
			Desc:   "Missing Username",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: "",
				Password: "TestPasswd@nsaop2021",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFieldEmpty,
		},
		{
			Desc:   "Missing Password",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: test.CommonUsers[constant.RoleCustomer].Username,
				Password: "",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFieldEmpty,
		},
		{
			Desc:   "Wrong Username",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: "NSAOP",
				Password: "TestPasswd",
			},
			Code: http.StatusUnauthorized,
			Msg:  "user not found",
		},
		{
			Desc:   "Wrong Password",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: test.CommonUsers[constant.RoleCustomer].Username,
				Password: "TestPasswd@nsaop2022",
			},
			Code: http.StatusUnauthorized,
			Msg:  constant.MsgWrongPassword,
		},
		{
			Desc:   "username too long",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: "这是一个二十字的名字这是一个二十字的名字超了",
				Password: "TestPasswd@nsaop2021",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   "password too long",
			Url:    url,
			Method: http.MethodPost,
			Param: user.UserLogin{
				Username: "Username",
				Password: "TestPasswd@nsaop2021TestPasswd@nsaop2021TestPasswd@nsaop2021TestPasswd@nsaop2021TestPasswd@nsaop2021TestPasswd@nsaop2021",
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
