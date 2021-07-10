package location_test

import (
	"net/http"
	"nsaop/model"
	"nsaop/router/api/v2/location"
	"nsaop/utils"
	"nsaop/utils/constant"
	"nsaop/utils/test"
	"testing"
)

func TestLocationId(t *testing.T) {
	url := "/location/%v"
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
		otherToken, _ := test.LoginAndGetToken(test.CommonUsers["other"])

		badCases := []test.TestCase{
			{
				Desc:   constant.MsgFormatError,
				Url:    utils.AddIDToPath(url, "wrong"),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   constant.MsgFormatError,
				Url:    utils.AddIDToPath(url, "-1"),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   constant.MsgIdNotAvailable,
				Url:    utils.AddIDToPath(url, 1),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + otherToken,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgIdNotAvailable,
			},
			{
				Desc:   constant.MsgMethodNotAllowed,
				Url:    utils.AddIDToPath(url, 1),
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

func TestGetLocationId(t *testing.T) {
	url := "/location/%v"

	customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

	var locations []model.Location
	model.DB.Find(&locations)
	var ret []location.LocationInfo
	for i := 3; i < 3+4; i++ {
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
			Url:    utils.AddIDToPath(url, ret[1].ID),
			Method: http.MethodGet,
			Headers: map[string]string{
				"Authorization": "Bearer " + customerToken,
			},
			Code: http.StatusOK,
			Msg:  constant.MsgOK,
			Data: ret[1],
		},
	}

	t.Run("Succeed", func(t *testing.T) {
		t.Parallel()
		test.PerformTest(t, goodCases, false)
	})
}

func TestChangeLocationId(t *testing.T) {
	url := "/location/%v"

	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])
		otherToken, _ := test.LoginAndGetToken(test.CommonUsers["other"])

		var locations []model.Location
		model.DB.Find(&locations)
		var ret []location.LocationInfo
		for i := 3; i < 3+4; i++ {
			loc := locations[i]
			ret = append(ret, location.LocationInfo{
				ID:      loc.ID,
				Comment: loc.Comment,
				Address: loc.Address,
				Contact: loc.Contact,
				Phone:   loc.Phone,
			})
		}

		badCases := []test.TestCase{
			{
				Desc:   constant.MsgFormatError,
				Url:    utils.AddIDToPath(url, ret[1].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: struct {
					Comment int `json:"comment"`
				}{
					Comment: 134,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "too long comment",
				Url:    utils.AddIDToPath(url, ret[1].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: location.LocationChange{
					Comment: "六〇九的宿舍冰箱商店实在是太方便了",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long address",
				Url:    utils.AddIDToPath(url, ret[1].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: location.LocationChange{
					Address: "这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址这是一个一百字的地址超了",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long contact",
				Url:    utils.AddIDToPath(url, ret[1].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: location.LocationChange{
					Contact: "这是一个十个字的名字这是一个十个字的名字超了",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   "too long phone",
				Url:    utils.AddIDToPath(url, ret[1].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: location.LocationChange{
					Phone: "13416219802134162198021341621980213298098123098124098",
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgInvalidInfo,
			},
			{
				Desc:   constant.MsgIdNotAvailable,
				Url:    utils.AddIDToPath(url, ret[1].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + otherToken,
				},
				Param: location.LocationChange{},
				Code:  http.StatusBadRequest,
				Msg:   constant.MsgIdNotAvailable,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "Change comment",
				Url:    utils.AddIDToPath(url, ret[0].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: location.LocationChange{
					Comment: "comment",
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Change others",
				Url:    utils.AddIDToPath(url, ret[0].ID),
				Method: http.MethodPut,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: location.LocationChange{
					Address: "华清大学",
					Contact: "联系人",
					Phone:   "98765432100",
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "check whether modify success",
				Url:    utils.AddIDToPath(url, ret[0].ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
				Data: location.LocationInfo{
					ID:      ret[0].ID,
					Comment: "comment",
					Address: "华清大学",
					Contact: "联系人",
					Phone:   "98765432100",
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
