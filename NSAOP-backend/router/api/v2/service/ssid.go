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

type AddSSIDObj struct {
	SSID SSIDObj `json:"ssid"`
}

type DelSSIDsObj struct {
	SSIDIDs []string `json:"ssid_ids"`
}

type SSIDObj struct {
	Name           string      `json:"name"`
	Enable         bool        `json:"enable"`
	ConnectionMode string      `json:"connection_mode"`
	HidedEnable    bool        `json:"hided_enable"`
	RelativeRadios int         `json:"relative_radios"`
	MaxUserNumber  int         `json:"max_user_number"`
	UserSeparation bool        `json:"user_separation"`
	SSIDAuth       SSIDAuthObj `json:"ssid_auth"`   // by default below
	SSIDPolicy     struct{}    `json:"ssid_policy"` // {}
}

type SSIDAuthObj struct {
	Mode   string    `json:"mode"`   // "open"
	Portal PortalObj `json:"portal"` // "portalDisable"
}

type PortalObj struct {
	Mode string `json:"mode"`
}

var (
	defaultPortal   = PortalObj{Mode: "portalDisable"}
	defaultSSIDAuth = SSIDAuthObj{
		Mode:   "open",
		Portal: defaultPortal,
	}
)

func SSID(c *gin.Context) {
	if c.Request.Method != http.MethodGet && c.MustGet("userRole") != constant.RoleEngineer {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}
	switch c.Request.Method {
	case http.MethodPost:
		AddSSID(c)
	case http.MethodGet:
		GetSSIDs(c)
	case http.MethodDelete:
		DelSSIDs(c)
	default:
		resp.ERROR(c, http.StatusMethodNotAllowed, constant.MsgMethodNotAllowed)
	}
}

// @Tags		Service
// @Summary		Add SSID when deploying service
// @Description	with jwt auth (by engineer), add SSID to a certain service
// @Description **must fill:**
// @Description <pre>{
// @Description		ssid, // object of a SSID (not array)
// @Description }</pre>
// @Description Object of ssid:
// @Description <pre>{
// @Description name, enable (bool), connection_mode (bridge / nat), hided_enable (bool),
// @Description relative_radios (int, 1 ~ 7), max_user_number (int, 1 ~ 512), user_separation (bool)
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param		Authorization header string true "Bearer [access_token]"
// @Param 		body body AddSSIDObj true "service's uuid + ssid"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>id not available"
// @Router		/service/{id}/ssid [POST]
func AddSSID(c *gin.Context) {
	var addSSID AddSSIDObj
	if err := c.ShouldBindJSON(&addSSID); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	service := c.MustGet(constant.TableService).(model.Service)
	if service.Status != constant.StatusOn && service.Status != constant.StatusPause && service.Status != constant.StatusPass {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}
	SSIDPath := fmt.Sprintf(constant.SsidFormatUrl, service.NCESiteID)
	// create SSID
	if len(addSSID.SSID.Name) == 0 {
		resp.OK(c, http.StatusOK, gin.H{})
		return
	}
	crtSSIDRqstBody := gin.H{
		"name":           addSSID.SSID.Name,
		"enable":         addSSID.SSID.Enable,
		"connectionMode": addSSID.SSID.ConnectionMode,
		"hidedEnable":    addSSID.SSID.HidedEnable,
		"relativeRadios": addSSID.SSID.RelativeRadios,
		"maxUserNumber":  addSSID.SSID.MaxUserNumber,
		"userSeparation": addSSID.SSID.UserSeparation,
		"ssidAuth":       defaultSSIDAuth,
		"ssidPolicy":     struct{}{},
	}
	// perform request to NCE
	crtSSIDResp := utils.RequestNCEWithToken(http.MethodPost, SSIDPath, nil, nil, crtSSIDRqstBody)
	if crtSSIDResp.Err != nil {
		resp.ERROR(c, crtSSIDResp.StatusCode, crtSSIDResp.Err.Error())
		return
	} else if crtSSIDResp.StatusCode != http.StatusOK &&
		crtSSIDResp.StatusCode != http.StatusCreated {
		resp.JSON(c, crtSSIDResp.StatusCode,
			crtSSIDResp.RespBody,
			fmt.Sprintf(constant.HttpCodeErrorFormat, crtSSIDResp.StatusCode))
		return
	} else {
		resp.OK(c, http.StatusOK, crtSSIDResp.RespBody)
		return
	}
}

// @Tags		Service
// @Summary		Get SSID of a service
// @Description	with jwt auth, get SSID of a certain service
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>id not available<br>"
// @Router		/service/{id}/ssid [GET]
func GetSSIDs(c *gin.Context) {
	service := c.MustGet(constant.TableService).(model.Service)
	if service.Status != constant.StatusOn && service.Status != constant.StatusPause && service.Status != constant.StatusPass && service.Status != constant.StatusSuspend && service.Status != constant.StatusRetrieve {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}
	SSIDPath := fmt.Sprintf(constant.SsidFormatUrl, service.NCESiteID)
	// get all SSIDs of the site
	getSSIDsResp := utils.RequestNCEWithToken(
		http.MethodGet, SSIDPath, nil, nil, nil)
	if getSSIDsResp.Err != nil {
		resp.ERROR(c, getSSIDsResp.StatusCode, getSSIDsResp.Err.Error())
		return
	} else if getSSIDsResp.StatusCode != http.StatusOK {
		resp.JSON(c, getSSIDsResp.StatusCode,
			getSSIDsResp.RespBody,
			fmt.Sprintf(constant.HttpCodeErrorFormat, getSSIDsResp.StatusCode))
		return
	} else {
		resp.OK(c, http.StatusOK, getSSIDsResp.RespBody)
		return
	}
}

// @Tags		Service
// @Summary		Delete SSIDs by given IDs
// @Description	with jwt auth (by engineer), delete SSIDs of a certain service
// @Description **must fill:**
// @Description <pre>{
// @Description		ssid_ids, // array of ids of the SSIDs that belong to the service
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param 		id path string ture "service id"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		body body DelSSIDsObj true "service's uuid"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>id not available<br>permission denied"
// @Router		/service/{id}/ssid [DELETE]
func DelSSIDs(c *gin.Context) {
	var delSSIDs DelSSIDsObj
	if err := c.ShouldBindJSON(&delSSIDs); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	service := c.MustGet(constant.TableService).(model.Service)
	SSIDPath := fmt.Sprintf(constant.SsidFormatUrl, service.NCESiteID)
	if len(delSSIDs.SSIDIDs) == 0 {
		resp.OK(c, http.StatusOK, gin.H{})
		return
	}
	// perform request to NCE
	delSSIDsRqstBody := gin.H{
		"ids": delSSIDs.SSIDIDs,
	}
	delSSIDsResp := utils.RequestNCEWithToken(
		http.MethodDelete, SSIDPath, nil, nil, delSSIDsRqstBody)
	if delSSIDsResp.Err != nil {
		resp.ERROR(c, delSSIDsResp.StatusCode, delSSIDsResp.Err.Error())
		return
	} else if delSSIDsResp.StatusCode != http.StatusOK {
		resp.JSON(c, delSSIDsResp.StatusCode,
			delSSIDsResp.RespBody,
			fmt.Sprintf(constant.HttpCodeErrorFormat, delSSIDsResp.StatusCode))
		return
	} else {
		resp.OK(c, http.StatusOK, delSSIDsResp.RespBody)
		return
	}
}
