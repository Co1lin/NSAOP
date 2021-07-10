package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"nsaop/model"
	"nsaop/router/resp"
	"nsaop/utils"
	"nsaop/utils/constant"
)

type AddDevicesObj struct {
	Devices []Device `json:"devices"`
}

type DelDevicesObj struct {
	DeviceIDs []string `json:"device_ids"`
}

type Device struct {
	Name        string `json:"name"`
	DeviceModel string `json:"device_model"`
}

var NCEDevicePath = "/controller/campus/v3/devices"

func Devices(c *gin.Context) {
	if c.Request.Method != http.MethodGet && c.MustGet("userRole") != constant.RoleEngineer {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}
	switch c.Request.Method {
	case http.MethodPost:
		AddDevices(c)
	case http.MethodGet:
		GetDevices(c)
	case http.MethodDelete:
		DelDevices(c)
	default:
		resp.ERROR(c, http.StatusMethodNotAllowed, "")
	}
}

// @Tags		Service
// @Summary		Add devices when deploying service
// @Description	with jwt auth (by engineer), add devices to a certain service
// @Description **must fill:**
// @Description <pre>{
// @Description		devices, // array of Device
// @Description }</pre>
// @Description *can be null:*
// @Description <pre>{
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		body body AddDevicesObj true "service's uuid + devices"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>permission denied<br>id not available<br>"
// @Router		/service/{id}/device [POST]
func AddDevices(c *gin.Context) {
	var addDevices AddDevicesObj
	if err := c.ShouldBindJSON(&addDevices); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	service := c.MustGet(constant.TableService).(model.Service)
	if service.Status != constant.StatusOn && service.Status != constant.StatusPause && service.Status != constant.StatusPass {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}
	// create devices
	crtDevRqstBody := gin.H{
		"devices": []gin.H{},
	}
	if len(addDevices.Devices) == 0 {
		resp.OK(c, http.StatusOK, gin.H{})
		return
	}
	for _, v := range addDevices.Devices {
		deviceInfo := gin.H{
			"name":        v.Name,
			"deviceModel": v.DeviceModel,
			"siteId":      service.NCESiteID,
		}
		crtDevRqstBody["devices"] = append(crtDevRqstBody["devices"].([]gin.H), deviceInfo)
	}
	// perform request to NCE
	crtDevResp := utils.RequestNCEWithToken(
		http.MethodPost, NCEDevicePath, nil, nil, crtDevRqstBody)
	if crtDevResp.Err != nil {
		resp.ERROR(c, crtDevResp.StatusCode, crtDevResp.Err.Error())
		return
	} else if crtDevResp.StatusCode != http.StatusOK {
		resp.JSON(c, crtDevResp.StatusCode,
			crtDevResp.RespBody,
			fmt.Sprintf(constant.HttpCodeErrorFormat, crtDevResp.StatusCode))
		return
	} else {
		resp.OK(c, http.StatusOK, crtDevResp.RespBody)
		return
	}
}

// @Tags		Service
// @Summary		Get devices of a service
// @Description	with jwt auth, get devices of a certain service
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "id not available<br>"
// @Router		/service/{id}/device [GET]
func GetDevices(c *gin.Context) {
	service := c.MustGet(constant.TableService).(model.Service)
	if service.NCESiteID == "" {
		resp.OK(c, http.StatusOK, gin.H{})
		return
	}
	// get all devices of the site
	getDevRqstParams := map[string]string{
		"siteId": service.NCESiteID,
	}
	// perform request to NCE
	getDevResp := utils.RequestNCEWithToken(
		http.MethodGet, NCEDevicePath, nil, getDevRqstParams, nil)
	if getDevResp.Err != nil {
		resp.ERROR(c, getDevResp.StatusCode, getDevResp.Err.Error())
		return
	} else if getDevResp.StatusCode != http.StatusOK {
		resp.JSON(c, getDevResp.StatusCode,
			getDevResp.RespBody,
			fmt.Sprintf(constant.HttpCodeErrorFormat, getDevResp.StatusCode))
		return
	} else {
		resp.OK(c, http.StatusOK, getDevResp.RespBody)
		return
	}
}

// @Tags		Service
// @Summary		Delete devices by given IDs
// @Description	with jwt auth (by engineer), delete devices of a certain service
// @Description **must fill:**
// @Description <pre>{
// @Description		device_ids, // array of ids of the devices
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		body body DelDevicesObj true "array of devices to delete"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>permission denied<br>invalid info<br>"
// @Router		/service/{id}/device [DELETE]
func DelDevices(c *gin.Context) {
	var delDevices DelDevicesObj
	if err := c.ShouldBindJSON(&delDevices); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if len(delDevices.DeviceIDs) == 0 {
		resp.OK(c, http.StatusOK, gin.H{})
		return
	}
	// perform request to NCE
	delDeviceRqstBody := gin.H{
		"deviceIds": delDevices.DeviceIDs,
	}
	delDeviceResp := utils.RequestNCEWithToken(
		http.MethodDelete, NCEDevicePath, nil, nil, delDeviceRqstBody)
	if delDeviceResp.Err != nil {
		resp.ERROR(c, delDeviceResp.StatusCode, delDeviceResp.Err.Error())
		return
	} else if delDeviceResp.StatusCode != http.StatusOK {
		resp.JSON(c, delDeviceResp.StatusCode,
			delDeviceResp.RespBody,
			fmt.Sprintf(constant.HttpCodeErrorFormat, delDeviceResp.StatusCode))
		return
	} else {
		resp.OK(c, http.StatusOK, delDeviceResp.RespBody)
		return
	}
}
