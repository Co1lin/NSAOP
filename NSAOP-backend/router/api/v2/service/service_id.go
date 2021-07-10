package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"nsaop/model"
	"nsaop/router/api/v2/location"
	"nsaop/router/resp"
	"nsaop/utils"
	"nsaop/utils/constant"
)

type ServiceChange struct {
	Target string `json:"target"`
	Stamp  uint   `json:"stamp"`
	Msg    string `json:"message"`
}

func ServiceId(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		GetServiceId(c)
	case http.MethodPut:
		ChangeServiceId(c)
	default:
		resp.ERROR(c, http.StatusMethodNotAllowed, constant.MsgMethodNotAllowed)
	}
}

// @Tags		Service
// @Summary		Get service info by id
// @Description	with jwt auth, get user's services detail info by id
// @Description	return
// @Description	<pre> {
// @Description   id, // uuid
// @Description   comment, // user's comment
// @Description   detail, // service detail
// @Description   paytype // = enum('month','year')
// @Description   require, // 0~7, 100 for private, 010 for client, 001 for test
// @Description   status, // = enum('waiting', 'pass', 'on', 'pause', 'canceled'); waiting： 等待运营师审核; pass： 审核完了等待工程师部署; on：正常工作中; pause：被暂停了，比如欠费; canceled：被删了
// @Description   create_at, // in format: "YYYY-MM-DD HH:mm:ss"
// @Description   pass_at, // in format: "YYYY-MM-DD HH:mm:ss"
// @Description   on_at, // in format: "YYYY-MM-DD HH:mm:ss"
// @Description   contact_operator, // operator's phone number
// @Description   contact_engineer, // engineer's phone number
// @Description   location_info, // same as detail in location's get api
// @Description	  message, // extra message
// @Description }<pre>
// @Accept		json
// @Product		json
// @Param 		id path string true "service id"
// @Param		Authorization header string true "Bearer [access_token]"
// @Success		200 {object} resp.Response{data=ServiceDetail} "ok"
// @Failure		400 {object} resp.Response "format error<br>id not available"
// @Router		/service/{id} [GET]
func GetServiceId(c *gin.Context) {
	service := c.MustGet(constant.TableService).(model.Service)
	var operator []model.User
	var engineer []model.User
	var loc model.Location
	model.DB.Model(&service).Where("role = ?", constant.RoleOperator).Association(constant.TableUser).Find(&operator)
	model.DB.Model(&service).Where("role = ?", constant.RoleEngineer).Association(constant.TableUser).Find(&engineer)
	model.DB.Unscoped().First(&loc, service.LocationID)
	resp.OK(c, http.StatusOK, ServiceDetail{
		ID:              service.ID.String(),
		Comment:         service.Comment,
		Detail:          service.Detail,
		PayType:         service.PayType,
		Require:         service.Require,
		Status:          service.Status,
		CreateAt:        model.Time2String(service.CreateAt),
		PassAt:          model.Time2String(service.PassAt),
		OnAt:            model.Time2String(service.OnAt),
		ContactOperator: operator[0].Phone,
		ContactEngineer: engineer[0].Phone,
		LocationInfo: location.LocationInfo{
			ID:      loc.ID,
			Comment: loc.Comment,
			Address: loc.Address,
			Contact: loc.Contact,
			Phone:   loc.Phone,
		},
		Msg: service.Msg,
	})
}

// @Tags		Service
// @Summary		Change service status
// @Description	with jwt auth, change user's services (by id) to given status
// @Description **must fill:**
// @Description <pre>{
// @Description		target, // target status
// @Description }</pre>
// @Description *can be null:*
// @Description <pre>{
// @Description		stamp, // needed if and only if operator want to approve
// @Description		message, // needed if and only if operator want to deny (1 <= len <= 200)
// @Description }</pre>
// @Accept		json
// @Product		json
// @Param 		id path string true "Bearer [access_token]"
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		body body ServiceChange true "service's uuid + target status"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "format error<br>invalid info<br>permission denied<br>id not available<br>not latest<br>undo limit exceed"
// @Router		/service/{id} [PUT]
func ChangeServiceId(c *gin.Context) {
	service := c.MustGet(constant.TableService).(model.Service)
	var change ServiceChange
	if err := c.ShouldBindJSON(&change); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if !model.ValidStatus(change.Target) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	if c.MustGet("userRole") == constant.RoleCustomer {
		switch source, target := service.Status, change.Target; {
		case source == constant.StatusWaiting && target == constant.StatusCanceled:
		case source == constant.StatusPass && target == constant.StatusCanceled:
		case source == constant.StatusOn && target == constant.StatusPause:
		case source == constant.StatusPause && target == constant.StatusOn:
		case source == constant.StatusPause && target == constant.StatusRetrieve:
		case source == constant.StatusRetrieve && target == constant.StatusPause:
			if service.Stamp >= 2 {
				resp.ERROR(c, http.StatusBadRequest, "undo limit exceed")
				return
			}
			service.Stamp += 1
		case source == constant.StatusSuspend && target == constant.StatusOn:
		default:
			resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
			return
		}
	} else if c.MustGet("userRole") == constant.RoleOperator {
		switch source, target := service.Status, change.Target; {
		case source == constant.StatusWaiting && target == constant.StatusPass:
			if change.Stamp != service.Stamp {
				resp.ERROR(c, http.StatusBadRequest, "not latest")
				return
			}
			if NCESiteID, err := utils.CreateNCESite(service.ID.String()); err != nil {
				resp.ERROR(c, http.StatusInternalServerError, err.Error())
				return
			} else {
				service.NCESiteID = NCESiteID
			}
			service.Stamp = 0
			service.PassAt = time.Now()
		case source == constant.StatusWaiting && target == constant.StatusCanceled:
			if change.Msg == "" || !model.ValidDetailLen(change.Msg) {
				resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
				return
			}
			service.Msg = change.Msg
		case source == constant.StatusCanceled:
			resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
			return
		default:
		}
	} else if c.MustGet("userRole") == constant.RoleEngineer {
		switch source, target := service.Status, change.Target; {
		case source == constant.StatusPass && target == constant.StatusOn:
			service.OnAt = time.Now()
		case source == constant.StatusRetrieve && target == constant.StatusCanceled:
		default:
			resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
			return
		}
	}
	if change.Target == constant.StatusCanceled {
		if service.NCESiteID != "" {
			if err := utils.DeleteNCESites([]string{service.NCESiteID}); err != nil {
				resp.ERROR(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}
	service.Status = change.Target
	model.DB.Save(&service)
	resp.OK(c, http.StatusOK, gin.H{})
}
