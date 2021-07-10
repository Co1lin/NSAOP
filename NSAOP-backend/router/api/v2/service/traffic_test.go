package service_test

import (
	"net/http"
	"nsaop/utils"
	"nsaop/utils/constant"
	"testing"
	"time"

	"nsaop/model"
	"nsaop/router/api/v2/service"
	"nsaop/utils/test"
)

func TestGetTraffic(t *testing.T) {
	url := "/service/%v/traffic"
	// separate DB env
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		otherToken, _ := test.LoginAndGetToken(test.CommonUsers["other"])
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

		// get service
		var services []model.Service
		model.DB.Where("status = ?", constant.StatusOn).Find(&services)
		testService := services[0]
		testService.OnAt = time.Now()
		model.DB.Save(&testService)

		var passedService model.Service
		model.DB.Where("status = ?", constant.StatusPass).First(&passedService)

		badCases := []test.TestCase{
			{
				Desc:   "Not On",
				Url:    utils.AddIDToPath(url, passedService.ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Invalid User",
				Url:    utils.AddIDToPath(url, testService.ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + otherToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgIdNotAvailable,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "",
				Url:    utils.AddIDToPath(url, testService.ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
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

func TestGetTrafficByID(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		service.GetTrafficByID(time.Now(), time.Now())
	})
}
