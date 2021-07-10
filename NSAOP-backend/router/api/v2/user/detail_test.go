package user_test

import (
	"net/http"
	"testing"

	"nsaop/router/api/v2/user"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestGetDetail(t *testing.T) {
	url := "/user/detail"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	goodCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: user.UserDetail{
				Role:     test.CommonUsers[constant.RoleCustomer].Role,
				Username: test.CommonUsers[constant.RoleCustomer].Username,
				Company:  test.CommonUsers[constant.RoleCustomer].Company,
				Phone:    test.CommonUsers[constant.RoleCustomer].Phone,
				Email:    test.CommonUsers[constant.RoleCustomer].Email,
			},
		},
	}

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})

	badCases := []test.TestCase{
		{
			Desc:   "Token Error #1",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer" + customerToken,
			},
			Code: http.StatusUnauthorized,
			Msg:  "authHeader format error",
		},
		{
			Desc:   "Token Error #2",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " +
					customerToken[:len(customerToken)-1] + "?",
			},
			Code: http.StatusUnauthorized,
			Msg:  "Couldn't handle this token",
		},
		{
			Desc:   "Missing header",
			Url:    url,
			Method: http.MethodGet,
			Code:   http.StatusUnauthorized,
			Msg:    "authHeader not found",
		},
		{
			Desc:   "Wrong header",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"AutHoRiZaTi0n": "Bearer " + customerToken,
			},
			Code: http.StatusUnauthorized,
			Msg:  "authHeader not found",
		},
		{
			Desc:   constant.MsgMethodNotAllowed,
			Url:    url,
			Method: http.MethodDelete,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Code: http.StatusMethodNotAllowed,
			Msg:  constant.MsgMethodNotAllowed,
		},
	}

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, badCases, false)
	})
}

func TestModifyDetail(t *testing.T) {
	url := "/user/detail"

	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		token, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

		badCases := []test.TestCase{
			{
				Desc:   constant.MsgFormatError,
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: struct {
					Phone int `json:"phone"`
				}{
					Phone: 1,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "invalid email",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					Email: "abc@cba",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Too long Password",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					Password: "ModifiedCompanyModifiedCompanyModifiedCompanyModifiedCompanyModifiedCompany",
					Phone:    "192837465",
					Email:    "ModifiedEmail@qq.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Too long Company",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					Company: "这是一个三十字的公司这是一个三十字的公司这是一个三十字的公司超了",
					Phone:   "192837465",
					Email:   "ModifiedEmail@qq.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Too long Email",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					Company: "ModifiedCompany",
					Phone:   "192837465",
					Email:   "test@nsaop.com.nsaopp.nsaop.nsaop.nsaop.nsaop.nsaop.nsaop",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Too long Phone",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					Company: "ModifiedCompany",
					Phone:   "192837465192837465192837465192837465192837465192837465192837465",
					Email:   "mail@qq.com",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "Wrong password",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					OldPassword: "WrongPasswd@nsaop2021",
					Password:    "ModifiedPasswd@nsaop2021",
				},
				Code: http.StatusUnauthorized,
				Msg:  constant.MsgWrongPassword,
			},
			{
				Desc:   "Mismatch passwords info",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					OldPassword: "",
					Password:    "ModifiedPasswd@nsaop2021",
				},
				Code: http.StatusUnauthorized,
				Msg:  constant.MsgWrongPassword,
			},
			{
				Desc:   "Mismatch passwords info",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					OldPassword: ";dljfsadljflashf",
					Password:    "",
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "Change password",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					OldPassword: "TestPasswd@nsaop2021",
					Password:    "ModifiedPasswd@nsaop2021",
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Change others",
				Url:    url,
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: user.UserModify{
					Company: "ModifiedCompany",
					Phone:   "192837465",
					Email:   "ModifiedEmail@qq.com",
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "check whether modify success",
				Url:    url,
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
				Data: user.UserDetail{
					Role:     test.CommonUsers[constant.RoleCustomer].Role,
					Username: test.CommonUsers[constant.RoleCustomer].Username,
					Company:  "ModifiedCompany",
					Phone:    "192837465",
					Email:    "ModifiedEmail@qq.com",
				},
			},
		}

		for _, myCase := range goodCases {
			t.Run("Succeed", func(t *testing.T) {
				test.PerformTest(t, []test.TestCase{myCase}, false)
			})
		}
	})
}
