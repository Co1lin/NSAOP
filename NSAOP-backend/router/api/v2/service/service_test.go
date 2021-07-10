package service_test

import (
	"net/http"
	"nsaop/model"
	"nsaop/router/api/v2/service"
	"nsaop/utils/constant"
	"nsaop/utils/test"
	"testing"
)

func TestNewService(t *testing.T) {
	url := "/service"
	defer test.SetupDatabaseForTest([]string{constant.TableUser, constant.TableLocation})()
	t.Run("", func(t *testing.T) {

		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
		operatorToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleOperator])
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])

		badCases := []test.TestCase{
			{
				Desc:   "Wrong Method",
				Url:    url,
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusMethodNotAllowed,
				Msg:  constant.MsgMethodNotAllowed,
			},
			{
				Desc:   "wrong format",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: struct {
					Require string `json:"require"`
				}{
					Require: "wrong",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "empty comment",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "",
					Detail:  "detail",
					PayType: constant.PaytypeMonth,
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "empty detail",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "",
					PayType: constant.PaytypeMonth,
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "empty pay type",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: "",
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "too long comment",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "commentdsakjfksajdflkjdsalkfjdsalkfjlksajflksa",
					Detail:  "detail",
					PayType: constant.PaytypeMonth,
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long detail",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetaildetail",
					PayType: constant.PaytypeMonth,
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "wrong pay type",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: "wrong",
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "wrong require",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: constant.PaytypeYear,
					LocId:   2,
					Require: -1,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   constant.MsgPermissionDenied,
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + operatorToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: constant.PaytypeMonth,
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
			{
				Desc:   constant.MsgPermissionDenied,
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: constant.PaytypeMonth,
					LocId:   2,
					Require: 7,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
			{
				Desc:   "not exists location",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: constant.PaytypeMonth,
					LocId:   0,
					Require: 7,
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
				Desc:   "",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.ServiceNew{
					Comment: "comment",
					Detail:  "detail",
					PayType: constant.PaytypeMonth,
					LocId:   1,
					Require: 7,
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

func BenchmarkNewService(b *testing.B) {
	defer test.SetupDatabaseForTest([]string{constant.TableUser, constant.TableLocation})()
	for i := 0; i < 10000; i++ {
		for _, service := range test.CommonServices {
			model.DB.Create(service)
		}
	}

	url := "/service"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	testcase := test.TestCase{
		Desc:   "",
		Url:    url,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Authorization": "Bearer " + customerToken,
		},
		Param: service.ServiceNew{
			Comment: "这是一个十个字的备注",
			Detail:  "这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情这是一个两百字的详情",
			PayType: constant.PaytypeMonth,
			LocId:   1,
			Require: 7,
		},
		Code: http.StatusOK,
		Msg:  constant.MsgOK,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			test.RequestOnly(testcase)
		}
	})
}

func TestGetServiceByFilter(t *testing.T) {
	url := "/service"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
	operatorToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleOperator])
	engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])

	var services []model.Service
	model.DB.Where("status = 'on' or status = 'pause'").Order("create_at desc").Find(&services)
	var ret []service.ServiceAbstract
	for i := 2; i < 2+3; i++ {
		s := services[i]
		ret = append(ret, service.ServiceAbstract{
			ID:       s.ID.String(),
			Comment:  s.Comment,
			PayType:  s.PayType,
			Require:  s.Require,
			Status:   s.Status,
			CreateAt: model.Time2String(s.CreateAt),
			Msg:      s.Msg,
		})
	}
	goodCases := []test.TestCase{
		{
			Desc:   "no search",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: service.ServiceFilter{
				Offset: 2,
				Limit:  3,
				Status: []string{constant.StatusOn, constant.StatusPause},
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: service.ServiceArray{
				Count:    2 * 8,
				Services: ret,
			},
		},
	}
	model.DB.Where("(status = 'on' or status = 'pause' or status = 'retrieve') and (comment like '%nt1%')").Order("create_at desc").Find(&services)
	ret = []service.ServiceAbstract{}
	services = services[:3]
	for i := 0; i < 3; i++ {
		s := services[i]
		ret = append(ret, service.ServiceAbstract{
			ID:       s.ID.String(),
			Comment:  s.Comment,
			PayType:  s.PayType,
			Require:  s.Require,
			Status:   s.Status,
			CreateAt: model.Time2String(s.CreateAt),
		})
	}
	goodCases = append(goodCases, []test.TestCase{
		{
			Desc:   "use search",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: service.ServiceFilter{
				Limit:  3,
				Status: []string{constant.StatusOn, constant.StatusPause, constant.StatusRetrieve},
				Search: "nt1",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: service.ServiceArray{
				Count:    6,
				Services: ret,
			},
		},
		{
			Desc:   "use search",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: service.ServiceFilter{
				Limit:  3,
				Status: []string{constant.StatusOn, constant.StatusPause, constant.StatusRetrieve},
				Search: ret[0].ID,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: service.ServiceArray{
				Count:    1,
				Services: []service.ServiceAbstract{ret[0]},
			},
		},
	}...)

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})

	badCases := []test.TestCase{
		{
			Desc:   constant.MsgInvalidInfo,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + operatorToken,
			},
			Param: struct {
				Limit uint `json:"limit"`
			}{
				Limit: 0,
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   constant.MsgInvalidInfo,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + operatorToken,
			},
			Param: struct {
				Limit uint `json:"limit"`
			}{
				Limit: 0,
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   constant.MsgInvalidInfo,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + engineerToken,
			},
			Param: struct {
				Status []string `json:"status"`
			}{
				Status: []string{constant.StatusWaiting, constant.StatusPass, "wrong"},
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   constant.MsgInvalidInfo,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + engineerToken,
			},
			Param: struct {
				Status []string `json:"status"`
			}{
				Status: []string{constant.StatusWaiting, constant.StatusPass, constant.StatusWaiting},
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

func BenchmarkGetServiceByFilter(b *testing.B) {
	defer test.SetupDatabaseForTest([]string{constant.TableUser, constant.TableLocation})()
	for i := 0; i < 10000; i++ {
		for _, service := range test.CommonServices {
			model.DB.Create(service)
		}
	}

	url := "/service"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	testcase := test.TestCase{
		Desc:   "use search",
		Url:    url,
		Method: http.MethodGet,
		Headers: map[string]string{
			"Authorization": "Bearer " + customerToken,
		},
		Param: service.ServiceFilter{
			Limit:  3,
			Status: []string{constant.StatusOn, constant.StatusPause, constant.StatusRetrieve},
			Search: "nt1 nt2 nt3 nt4 nt5 a b c d e f g h i j k l m n o p q r s t u v w x y z",
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			test.RequestOnly(testcase)
		}
	})
}
