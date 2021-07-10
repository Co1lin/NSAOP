package service_test

import (
	"net/http"
	"testing"

	"nsaop/model"
	"nsaop/utils"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestIdMiddleware(t *testing.T) {
	url := "/service/%v"
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
		otherToken, _ := test.LoginAndGetToken(test.CommonUsers["other"])

		var services []model.Service
		model.DB.Where("status = 'waiting' or status = 'pass'").Find(&services)

		badCases := []test.TestCase{
			{
				Desc:   constant.MsgFormatError,
				Url:    utils.AddIDToPath(url, services[0].ID.String()[:len(services[0].ID.String())-1]+"?"),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "Too long ID",
				Url:    utils.AddIDToPath(url, services[1].ID.String()+"long------------------------------------------------------------------"),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   constant.MsgIdNotAvailable,
				Url:    utils.AddIDToPath(url, services[1].ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + otherToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgIdNotAvailable,
			},
			{
				Desc:   constant.MsgMethodNotAllowed,
				Url:    utils.AddIDToPath(url, services[0].ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusMethodNotAllowed,
				Msg:  constant.MsgMethodNotAllowed,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			test.PerformTest(t, badCases, false)
		})
	})
}
