package location_test

import (
	"net/http"
	"nsaop/model"
	"nsaop/router/api/v2/location"
	"nsaop/utils/constant"
	"nsaop/utils/test"
	"testing"
)

func TestNewLocation(t *testing.T) {
	url := "/location"

	defer test.SetupDatabaseForTest([]string{constant.TableUser})()
	t.Run("", func(t *testing.T) {
		token, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
		operatorToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleOperator])
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])

		badCases := []test.TestCase{
			{
				Desc:   "Wrong Method",
				Url:    url,
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Code: http.StatusMethodNotAllowed,
				Msg:  constant.MsgMethodNotAllowed,
			},
			{
				Desc:   "wrong format",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: struct {
					Comment int `json:"comment"`
				}{
					Comment: 1,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "empty comment",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "",
					Address: "紫荆公寓609B",
					Contact: "彭先生",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "empty address",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "",
					Contact: "彭先生",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "empty contact",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "紫荆公寓609B",
					Contact: "",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   "empty phone",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "紫荆公寓609B",
					Contact: "彭先生",
					Phone:   "",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFieldEmpty,
			},
			{
				Desc:   constant.MsgPermissionDenied,
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + operatorToken,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "紫荆公寓609B",
					Contact: "彭先生",
					Phone:   "13416219802",
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
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "紫荆公寓609B",
					Contact: "彭先生",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
			{
				Desc:   "too long comment",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "六〇九的宿舍冰箱商店实在是太方便了",
					Address: "紫荆公寓609B",
					Contact: "彭先生",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long address",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址超了",
					Contact: "彭先生",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long contact",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "紫荆公寓609B",
					Contact: "这是一个十个字的名字这是一个十个字的名字超了",
					Phone:   "13416219802",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long phone",
				Url:    url,
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "紫荆公寓609B",
					Contact: "彭先生",
					Phone:   "13416219802134162198021341621980213298098123098124098",
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
					"Authorization": "Bearer " + token,
				},
				Param: location.LocationNew{
					Comment: "宿舍冰箱商店",
					Address: "北京市海淀区清华园街道清华大学紫荆学生公寓609B",
					Contact: "彭先生",
					Phone:   "13416219802",
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

func TestGetLocationByFilter(t *testing.T) {
	url := "/location"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	var locations []model.Location
	model.DB.Find(&locations)
	var ret []location.LocationInfo
	for i := 2; i < 2+3; i++ {
		loc := locations[i]
		ret = append(ret, location.LocationInfo{
			ID:      loc.ID,
			Comment: loc.Comment,
			Address: loc.Address,
			Contact: loc.Contact,
			Phone:   loc.Phone,
		})
	}
	goodCases := []test.TestCase{
		{
			Desc:   "",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Offset: 2,
				Limit:  3,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: location.LocationArray{
				Locations: ret,
				Count:     len(test.CommonLocation),
			},
		},
		{
			Desc:   "search comment",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Limit:  1,
				Search: ret[0].Comment,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: location.LocationArray{
				Count:     11,
				Locations: []location.LocationInfo{ret[0]},
			},
		},
		{
			Desc:   "search address",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Limit:  1,
				Search: "63B",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: location.LocationArray{
				Count:     1,
				Locations: []location.LocationInfo{ret[0]},
			},
		},
		{
			Desc:   "search with special token _",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Limit:  1,
				Search: "63_",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: location.LocationArray{
				Count:     0,
				Locations: nil,
			},
		},
		{
			Desc:   "search with special token %",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Limit:  1,
				Search: "63%",
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: location.LocationArray{
				Count:     0,
				Locations: nil,
			},
		},
		{
			Desc:   `search with special token \%`,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Limit:  1,
				Search: `6\%`,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: location.LocationArray{
				Count:     0,
				Locations: nil,
			},
		},
	}

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})

	badCases := []test.TestCase{
		{
			Desc:   constant.MsgFormatError,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: struct {
				Offset string `json:"offset"`
			}{
				Offset: "wrong",
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgFormatError,
		},
		{
			Desc:   constant.MsgInvalidInfo,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Offset: 0,
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   constant.MsgInvalidInfo,
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Offset: 1000,
			},
			Code: http.StatusBadRequest,
			Msg:  constant.MsgInvalidInfo,
		},
		{
			Desc:   "search too long",
			Url:    url,
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Param: location.LocationFilter{
				Search: "tooooooooooooooooooooooooooooooo loooooooooooooooooooooooooooooooooooooooooooooooong loooooooooooooooooooooooooooooooooooooooooooooooooooong",
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
