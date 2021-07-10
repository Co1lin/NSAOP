package service_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"nsaop/model"
	"nsaop/router/api/v2/location"
	"nsaop/router/api/v2/service"
	"nsaop/utils"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

func TestGetServiceId(t *testing.T) {
	assert := assert.New(t)
	url := "/service/%v"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
	operatorToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleOperator])
	engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])

	var services []model.Service
	model.DB.Where("status = 'waiting' or status = 'pass'").Find(&services)
	var ret []service.ServiceDetail
	services = services[3 : 3+4]
	for i := 0; i < 4; i++ {
		s := services[i]
		var operator []model.User
		var engineer []model.User
		var loc model.Location
		model.DB.Model(&s).Where("role = ?", constant.RoleOperator).Association("Users").Find(&operator)
		model.DB.Model(&s).Where("role = ?", constant.RoleEngineer).Association("Users").Find(&engineer)
		model.DB.First(&loc, s.LocationID)
		assert.Equal(len(operator), 1, "one service has multiple operator")
		assert.Equal(len(engineer), 1, "one service has multiple engineer")
		ret = append(ret, service.ServiceDetail{
			ID:              s.ID.String(),
			Comment:         s.Comment,
			Detail:          s.Detail,
			PayType:         s.PayType,
			Require:         s.Require,
			Status:          s.Status,
			CreateAt:        model.Time2String(s.CreateAt),
			PassAt:          model.Time2String(s.PassAt),
			OnAt:            model.Time2String(s.OnAt),
			ContactOperator: operator[0].Phone,
			ContactEngineer: engineer[0].Phone,
			LocationInfo: location.LocationInfo{
				ID:      loc.ID,
				Comment: loc.Comment,
				Address: loc.Address,
				Contact: loc.Contact,
				Phone:   loc.Phone,
			},
			Msg: s.Msg,
		})
	}
	model.DB.Where(constant.IdQuery, services[1].LocationID).Delete(&model.Location{})
	goodCases := []test.TestCase{
		{
			Desc:   "customer get",
			Url:    utils.AddIDToPath(url, services[1].ID),
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: ret[1],
		},
		{
			Desc:   "operator get",
			Url:    utils.AddIDToPath(url, services[2].ID),
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + operatorToken,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: ret[2],
		},
		{
			Desc:   "engineer get",
			Url:    utils.AddIDToPath(url, services[0].ID),
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + engineerToken,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: ret[0],
		},
	}

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})
}

func TestChangeServiceId(t *testing.T) {
	url := "/service/%v"

	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
		operatorToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleOperator])
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])

		var services []model.Service
		model.DB.Where("status = 'waiting'").Find(&services)
		var ret []service.ServiceAbstract
		services = services[3 : 3+5]
		for i := 0; i < 5; i++ {
			s := services[i]
			ret = append(ret, service.ServiceAbstract{
				ID:       s.ID.String(),
				PayType:  s.PayType,
				Require:  s.Require,
				Status:   constant.StatusPass,
				CreateAt: model.Time2String(s.CreateAt),
			})
		}

		normalCases := []test.TestCase{
			{
				Desc:   "success case",
				Url:    utils.AddIDToPath(url, services[0].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + operatorToken,
				},
				Param: service.ServiceChange{
					Target: constant.StatusPass,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   constant.MsgInvalidInfo,
				Url:    utils.AddIDToPath(url, services[0].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceChange{
					Target: "wrong",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
		}

		t.Run("normal", func(t *testing.T) {
			test.PerformTest(t, normalCases, false)
		})

		trans := func(Desc string, token string, which int, Target string, Stamp uint, Message string) test.TestCase {
			Code := http.StatusOK
			Msg := constant.MsgOK
			if Desc != constant.MsgOK {
				Code = http.StatusBadRequest
				Msg = Desc
			}
			return test.TestCase{
				Desc:   Desc,
				Url:    utils.AddIDToPath(url, services[which].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: service.ServiceChange{
					Target: Target,
					Stamp:  Stamp,
					Msg:    Message,
				},
				Code: Code,
				Msg:  Msg,
			}
		}
		lifecycleCases := []test.TestCase{
			trans(constant.MsgPermissionDenied, customerToken, 1, constant.StatusPass, 0, ""),
			trans(constant.MsgPermissionDenied, customerToken, 1, constant.StatusOn, 0, ""),
			trans("not latest", operatorToken, 1, constant.StatusPass, 100, ""),
			trans(constant.MsgOK, operatorToken, 1, constant.StatusPass, 0, ""),
			trans(constant.MsgOK, engineerToken, 1, constant.StatusOn, 0, ""),
			trans(constant.MsgOK, operatorToken, 1, constant.StatusSuspend, 0, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusOn, 0, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusPause, 0, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusOn, 0, ""),
			trans(constant.MsgOK, operatorToken, 1, constant.StatusPause, 0, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusRetrieve, 0, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusPause, 3, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusRetrieve, 0, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusPause, 3, ""),
			trans(constant.MsgOK, customerToken, 1, constant.StatusRetrieve, 0, ""),
			trans("undo limit exceed", customerToken, 1, constant.StatusPause, 0, ""),
			trans(constant.MsgOK, engineerToken, 1, constant.StatusCanceled, 0, ""),
			trans(constant.MsgPermissionDenied, engineerToken, 1, constant.StatusRetrieve, 0, ""),
			trans(constant.MsgPermissionDenied, operatorToken, 1, constant.StatusOn, 0, ""),

			trans(constant.MsgOK, customerToken, 2, constant.StatusCanceled, 0, ""),
			trans(constant.MsgPermissionDenied, customerToken, 2, constant.StatusOn, 0, ""),

			trans(constant.MsgOK, operatorToken, 3, constant.StatusPass, 0, ""),

			trans(constant.MsgInvalidInfo, operatorToken, 4, constant.StatusCanceled, 0, ""),
			trans(constant.MsgInvalidInfo, operatorToken, 4, constant.StatusCanceled, 0, "10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			trans(constant.MsgOK, operatorToken, 4, constant.StatusCanceled, 0, "deny reason"),
		}

		for _, myCase := range lifecycleCases {
			t.Run("lifecycle", func(t *testing.T) {
				test.PerformTest(t, []test.TestCase{myCase}, false)
			})
		}
	})
}
