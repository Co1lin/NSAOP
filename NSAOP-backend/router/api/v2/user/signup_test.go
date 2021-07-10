package user_test

import (
	"net/http"
	"nsaop/config"
	"testing"

	"nsaop/router/api/v2/user"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestSignup(t *testing.T) {
	url := "/user/signup"
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		badCases := []test.TestCase{
			{
				Desc:   "Format error",
				Url:    url,
				Method: http.MethodPost,
				Param: struct {
					Username int `json:"username"`
				}{
					Username: 123,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "Empty Role",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Username: "EmptyRoleTest",
					Password: "MyPasswd",
					Company:  "Go Test",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "Invalid Role",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     "administrator",
					Username: "InvalidRoleTest",
					Password: "MyPasswd",
					Company:  "Go Test",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Conflicting Username",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: test.CommonUsers[constant.RoleCustomer].Username,
					Password: "MyPasswd",
					Company:  "Go Test",
				},
				Code: http.StatusBadRequest,
				Msg:  "username already exist",
			},
			{
				Desc:   "Missing Username",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Password: "MyPasswd",
					Company:  "Go Test",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "Missing Password",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestForMissingPasswd",
					Company:  "Go Test",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "Missing Name",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestForMissingName",
					Company:  "Go Test",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "Email Format Error #1",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestForEmailError",
					Company:  "Go Test",
					Password: "MyPasswd",
					Email:    "testnsaop",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Email Format Error #2",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestForEmailError",
					Company:  "Go Test",
					Password: "MyPasswd",
					Phone:    "98457834658",
					Email:    "test@nsaop",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Email Format Error #3",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestForEmailError",
					Company:  "Go Test",
					Password: "MyPasswd",
					Phone:    "98457834658",
					Email:    "nsaop_email@",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Username Format Error (len > 20)",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleOperator,
					Username: "这是一个二十字的名字这是一个二十字的名字超了",
					Company:  "Go Test",
					Password: "MyPasswd",
					Phone:    "98457834658",
					Email:    "email@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Username Format Error (len < 5)",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleEngineer,
					Username: "haha",
					Company:  "Go Test",
					Password: "MyPasswd",
					Phone:    "98457834658",
					Email:    "email@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Too long",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestCustomerGoTestCustomerGoTestCustomerGoTestCustomerGoTestCustomerGoTestCustomer",
					Password: "MyPasswd",
					Company:  "Go Test",
					Phone:    "98457834658",
					Email:    "test@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Too long",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestCustomer",
					Password: "MyPasswdMyPasswdMyPasswdMyPasswdMyPasswdMyPasswdMyPasswdMyPasswdMyPasswdMyPasswd",
					Company:  "Go Test",
					Phone:    "98457834658",
					Email:    "test@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "company name too long",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestCustomer",
					Password: "MyPasswd",
					Company:  "这是一个三十字的公司这是一个三十字的公司这是一个三十字的公司超了",
					Phone:    "98457834658",
					Email:    "test@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "phone too long",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestCustomer",
					Password: "MyPasswd",
					Company:  "Go Test",
					Phone:    "9845783465898457834658984578346589845783465898457834658",
					Email:    "test@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "email too long",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestCustomer",
					Password: "MyPasswd",
					Company:  "Go Test",
					Phone:    "98457834658",
					Email:    "test@nsaop.com.nsaopp.nsaop.nsaop.nsaop.nsaop.nsaop.nsaop",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Operator without register code",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleOperator,
					Username: "GoTestOperator",
					Password: "MyPasswd",
					Company:  "NSAOP",
					Phone:    "13357565444",
					Email:    "test@nsaop.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "Customer",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleCustomer,
					Username: "GoTestCustomer",
					Password: "MyPasswd",
					Company:  "Go Test",
					Phone:    "98457834658",
					Email:    "test@nsaop.com",
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Operator",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleOperator,
					Username: "GoTestOperator",
					Password: "MyPasswd",
					Company:  "NSAOP",
					Phone:    "13357565444",
					Email:    "test@nsaop.com",
					Code:     config.AdminRegisterCode,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Engineer",
				Url:    url,
				Method: http.MethodPost,
				Param: user.UserSignup{
					Role:     constant.RoleEngineer,
					Username: "GoTestEngineer",
					Password: "MyPasswd",
					Company:  "NSAOP",
					Phone:    "13357565444",
					Email:    "test@nsaop.com",
					Code:     config.AdminRegisterCode,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
		}
		t.Run("Succeed", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, goodCases, false)
		})
	})
}
