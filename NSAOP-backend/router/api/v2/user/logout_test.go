package user_test

import (
	"net/http"
	"testing"

	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestLogout(t *testing.T) {
	url := "/user/logout"

	_, refreshToken := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	badCases := []test.TestCase{
		{
			Desc:   "Wrong refresh_token value",
			Url:    url,
			Method: http.MethodPost,
			Cookies: []http.Cookie{
				{
					Name:  "refresh_token",
					Value: refreshToken[:len(refreshToken)-1],
				},
			},
			Code: http.StatusBadRequest,
			Msg:  "token not exists",
		},
		{
			Desc:   "Wrong refresh_token name",
			Url:    url,
			Method: http.MethodPost,
			Cookies: []http.Cookie{
				{
					Name:  "refresh token",
					Value: refreshToken,
				},
			},
			Code: http.StatusBadRequest,
			Msg:  "refresh_token not found",
		},
		{
			Desc:    "No Cookies",
			Url:     url,
			Method:  http.MethodPost,
			Cookies: []http.Cookie{},
			Code:    http.StatusBadRequest,
			Msg:     "refresh_token not found",
		},
	}
	// shouldn't be parallel since after a successful logout,
	// the cookies that have got will be invalid,
	// which causes the following tests meaningless
	t.Run("Fail", func(t *testing.T) {
		test.PerformTest(t, badCases, false)
	})

	goodCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodPost,
			Cookies: []http.Cookie{
				{
					Name:  "refresh_token",
					Value: refreshToken,
				},
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
		},
	}
	// after all failed logout test,
	// we can finally logout successfully
	t.Run("Succeed", func(t *testing.T) {
		test.PerformTest(t, goodCases, false)
	})
}
