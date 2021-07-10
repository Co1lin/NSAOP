package service_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"nsaop/model"
	"nsaop/router/api/v2/service"
	"nsaop/utils"
	"nsaop/utils/constant"
	"nsaop/utils/test"
)

const devicePath = "/service/%v/device"

func TestAddDevices(t *testing.T) {
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
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.AddDevicesObj{
					Devices: []service.Device{},
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgPermissionDenied,
			},
			{
				Desc:   "Invalid Status",
				Url:    utils.AddIDToPath(devicePath, testWrongStatusService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddDevicesObj{
					Devices: []service.Device{},
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
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddDevicesObj{
					Devices: []service.Device{
						{
							Name:        "virtual-test-device-for-add-01",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-add-02",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-add-03",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-add-04",
							DeviceModel: "AP4050DN",
						},
					},
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Empty",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddDevicesObj{
					Devices: []service.Device{},
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
	//DeleteDevicesOfTestSite()
}

func TestGetDevices(t *testing.T) {
	// separate DB env
	defer test.SetupDatabaseForTest(deferBase)()
	t.Run("", func(t *testing.T) {
		engineerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleEngineer])
		customerToken, _ := test.LoginAndGetToken(test.CommonUsers[constant.RoleCustomer])

		// get service
		var services []model.Service
		model.DB.Where("status = 'pass'").Find(&services)
		// use a different one to avoid collision with test function(s) above
		testService := services[1]
		testService.NCESiteID = test.CommonNCESiteID
		model.DB.Save(&testService)

		// create devices first
		caseForCreate := []test.TestCase{
			{
				Desc:   "Engineer",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddDevicesObj{
					Devices: []service.Device{
						{
							Name:        "virtual-test-device-for-get-01",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-get-02",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-get-03",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-get-04",
							DeviceModel: "AP4050DN",
						},
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
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Engineer",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
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
	//DeleteDevicesOfTestSite()
}

func TestDelDevices(t *testing.T) {
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

		// create devices first
		caseForCreate := []test.TestCase{
			{
				Desc:   "Engineer",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.AddDevicesObj{
					Devices: []service.Device{
						{
							Name:        "virtual-test-device-for-del-01",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-del-02",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-del-03",
							DeviceModel: "AP4050DN",
						},
						{
							Name:        "virtual-test-device-for-del-04",
							DeviceModel: "AP4050DN",
						},
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
				Desc:   "Format Error",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: struct {
					DeviceIDs int `json:"device_ids"`
				}{
					DeviceIDs: 1,
				},
				Code: http.StatusBadRequest,
				Msg:  constant.MsgFormatError,
			},
			{
				Desc:   "Invalid Role",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + customerToken,
				},
				Param: service.DelDevicesObj{
					DeviceIDs: []string{},
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
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.DelDevicesObj{
					DeviceIDs: getDevicesFromTestSite(),
				},
				Code: http.StatusOK,
				Msg:  constant.MsgOK,
			},
			{
				Desc:   "Empty",
				Url:    utils.AddIDToPath(devicePath, testService.ID),
				Method: http.MethodDelete,
				Headers: map[string]string{
					"Authorization": "Bearer " + engineerToken,
				},
				Param: service.DelDevicesObj{
					DeviceIDs: []string{},
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

func getDevicesFromTestSite() (devices []string) {
	getDevRqstParams := map[string]string{
		"siteId": test.CommonNCESiteID,
	}
	// get devices at the test site
	if devicesResp := utils.RequestNCEWithToken(
		http.MethodGet, service.NCEDevicePath, nil, getDevRqstParams, nil); devicesResp.Err != nil || devicesResp.StatusCode != http.StatusOK ||
		devicesResp.RespBody["data"] == nil ||
		len(devicesResp.RespBody["data"].([]interface{})) == 0 {
		//panic("Failed to get test devices")
		return
	} else {
		// get IDs of devices for test
		devicesList := devicesResp.RespBody["data"].([]interface{})
		for _, device := range devicesList {
			devices = append(devices, device.(map[string]interface{})["id"].(string))
		}
		return
	}
}

func DeleteDevicesOfTestSite() {
	devices := getDevicesFromTestSite()
	if devices == nil {
		return
	}
	delDeviceRqstBody := gin.H{
		"deviceIds": devices,
	}
	delDeviceResp := utils.RequestNCEWithToken(
		http.MethodDelete, service.NCEDevicePath, nil, nil, delDeviceRqstBody)
	if delDeviceResp.StatusCode != http.StatusOK || delDeviceResp.Err != nil {
		log.Fatalf("Failed to delete test devices; HTTP Code: %d\n%v\n", delDeviceResp.StatusCode, delDeviceResp.RespBody)
	} else {
		// log.Println("Test devices deleted.")
	}
}
