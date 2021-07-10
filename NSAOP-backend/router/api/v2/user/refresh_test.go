package user_test

import (
	"net/http"
	"testing"

	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestRefresh(t *testing.T) {
	url := "/user/refresh"
	logouturl := "/user/logout"

	_, refreshToken := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
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

	t.Run("Succeed", func(t *testing.T) {
		// cannot parralel, since logout in badCase
		test.PerformTest(t, goodCases, false)
	})

	badCases := []test.TestCase{
		{
			Desc:   "token not found",
			Url:    url,
			Method: http.MethodPost,
			Code:   http.StatusUnauthorized,
			Msg:    "token not found",
		},
		{
			Desc:   constant.MsgFormatError,
			Url:    url,
			Method: http.MethodPost,
			Cookies: []http.Cookie{
				{
					Name:  "refresh_token",
					Value: refreshToken[:len(refreshToken)-1],
				},
			},
			Code: http.StatusUnauthorized,
			Msg:  "Couldn't handle this token",
		},
		{
			Desc:   "logout",
			Url:    logouturl,
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
		{
			Desc:   "already logout",
			Url:    url,
			Method: http.MethodPost,
			Cookies: []http.Cookie{
				{
					Name:  "refresh_token",
					Value: refreshToken,
				},
			},
			Code: http.StatusUnauthorized,
			Msg:  "token has banned",
		},
	}

	for _, myCase := range badCases {
		t.Run("Fail", func(t *testing.T) {
			test.PerformTest(t, []test.TestCase{myCase}, false)
		})
	}
}
