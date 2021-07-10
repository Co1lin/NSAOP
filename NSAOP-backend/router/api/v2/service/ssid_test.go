package service_test

import (
	"log"
	"net/http"
	"nsaop/model"
	"nsaop/router/api/v2/service"
	"nsaop/utils"
	"nsaop/utils/constant"
	"nsaop/utils/test"
	"testing"

	"github.com/gin-gonic/gin"
)

const ssidPath = "/service/%v/ssid"

func TestAddSSID(t *testing.T) {
	// separate DB env
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

		// get service
		var services []model.Service
		model.DB.Where("status = 'pass'").Find(&services)
		testService := services[0]
		testService.NCESiteID = test.CommonNCESiteID
		model.DB.Save(&testService)

		testWrongStatusService := services[1]
		testWrongStatusService.Status = constant.StatusWaiting
		testWrongStatusService.NCESiteID = test.CommonNCESiteID
		model.DB.Save(&testWrongStatusService)

		badCases := []test.TestCase{
			{
				Desc:   "Invalid Role",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{},
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
			{
				Desc:   "Invalid Status",
				Url:    utils.AddIDToPath(ssidPath, testWrongStatusService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{},
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name:           "Guest-NSAOP-AddTest",
						Enable:         true,
						ConnectionMode: "bridge",
						HidedEnable:    true,
						RelativeRadios: 7,
						MaxUserNumber:  100,
						UserSeparation: true,
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Empty",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name: "",
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name:           "Employee-NSAOP-AddTest",
						Enable:         true,
						ConnectionMode: "bridge",
						HidedEnable:    true,
						RelativeRadios: 7,
						MaxUserNumber:  100,
						UserSeparation: true,
					},
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
	//DeleteSSIDsOfTestSite()
}

func TestGetSSIDs(t *testing.T) {
	// separate DB env
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

		// get service
		var services []model.Service
		model.DB.Where("status = 'pass'").Find(&services)
		// use a different one to avoid collision with test function(s) above
		testService := services[2]
		testService.NCESiteID = test.CommonNCESiteID
		model.DB.Save(&testService)

		// create SSIDs first
		caseForCreate := []test.TestCase{
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name:           "Guest-NSAOP-GetTest",
						Enable:         true,
						ConnectionMode: "bridge",
						HidedEnable:    true,
						RelativeRadios: 7,
						MaxUserNumber:  100,
						UserSeparation: true,
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name:           "Employee-NSAOP-GetTest",
						Enable:         true,
						ConnectionMode: "bridge",
						HidedEnable:    true,
						RelativeRadios: 7,
						MaxUserNumber:  100,
						UserSeparation: true,
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
		}

		t.Run("Create", func(t *testing.T) {
			test.PerformTest(t, caseForCreate, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "Customer",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Engineer",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
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
	//DeleteSSIDsOfTestSite()
}

func TestDelSSIDs(t *testing.T) {
	// separate DB env
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

		// get service
		var services []model.Service
		model.DB.Where("status = 'pass'").Find(&services)
		// use a different one to avoid collision with test function(s) above
		testService := services[3]
		testService.NCESiteID = test.CommonNCESiteID
		model.DB.Save(&testService)

		// create SSIDs first
		caseForCreate := []test.TestCase{
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name:           "Guest-NSAOP-DelTest",
						Enable:         true,
						ConnectionMode: "bridge",
						HidedEnable:    true,
						RelativeRadios: 7,
						MaxUserNumber:  100,
						UserSeparation: true,
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddSSIDObj{
					SSID: service.SSIDObj{
						Name:           "Employee-NSAOP-DelTest",
						Enable:         true,
						ConnectionMode: "bridge",
						HidedEnable:    true,
						RelativeRadios: 7,
						MaxUserNumber:  100,
						UserSeparation: true,
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
		}

		t.Run("Create", func(t *testing.T) {
			test.PerformTest(t, caseForCreate, false)
		})

		badCases := []test.TestCase{
			{
				Desc:   "Invalid Role",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.DelSSIDsObj{
					SSIDIDs: []string{},
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
		}

		t.Run("Fail", func(t *testing.T) {
			t.Parallel()
			test.PerformTest(t, badCases, false)
		})

		goodCases := []test.TestCase{
			{
				Desc:   "",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.DelSSIDsObj{
					SSIDIDs: getSSIDsFromTestSite(),
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Empty",
				Url:    utils.AddIDToPath(ssidPath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.DelSSIDsObj{
					SSIDIDs: []string{},
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

func getSSIDsFromTestSite() (ssids []string) {
	if ssidsResp := utils.RequestNCEWithToken(
		http.MethodGet, test.CommonNCESSIDPath, nil, nil, nil); ssidsResp.Err != nil || ssidsResp.StatusCode != http.StatusOK ||
		ssidsResp.RespBody["data"] == nil ||
		len(ssidsResp.RespBody["data"].([]interface{})) == 0 {
		//panic("Failed to get test SSIDs")
		return
	} else {
		// get IDs of SSIDs for test
		ssidsList := ssidsResp.RespBody["data"].([]interface{})
		for _, ssid := range ssidsList {
			ssids = append(ssids, ssid.(map[string]interface{})["id"].(string))
		}
		return
	}
}

func DeleteSSIDsOfTestSite() {
	ssids := getSSIDsFromTestSite()
	if ssids == nil {
		return
	}
	delSSIDsRqstBody := gin.H{
		"ids": ssids,
	}
	delSSIDsResp := utils.RequestNCEWithToken(
		http.MethodDelete, test.CommonNCESSIDPath, nil, nil, delSSIDsRqstBody)
	if delSSIDsResp.StatusCode != http.StatusOK || delSSIDsResp.Err != nil {
		log.Fatalf("Failed to delete test SSIDs; HTTP Code: %d\n%v\n", delSSIDsResp.StatusCode, delSSIDsResp.RespBody)
	} else {
		// log.Println("Test devices deleted.")
	}
}
