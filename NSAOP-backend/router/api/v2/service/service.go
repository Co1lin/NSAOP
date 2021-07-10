package service

import (
	"log"
	"math/rand"
	"net/http"
	"nsaop/model"
	"nsaop/router/api/v2/location"
	"nsaop/router/resp"
	"nsaop/utils"
	"nsaop/utils/constant"

	"github.com/gin-gonic/gin"
)

type ServiceFilter struct {
	Offset int      `json:"offset" form:"offset"`
	Limit  int      `json:"limit" form:"limit"`
	Status []string `json:"status" form:"status"`
	Search string   `json:"search" form:"search"`
}

type ServiceAbstract struct {
	ID       string `json:"id"`
	Comment  string `json:"comment"`
	PayType  string `json:"paytype"`
	Require  int    `json:"require"`
	Status   string `json:"status"`
	CreateAt string `json:"create_at"`
	Msg      string `json:"message"`
}

type ServiceArray struct {
	Count    int               `json:"count"`
	Services []ServiceAbstract `json:"services"`
}

type ServiceDetail struct {
	ID              string                `json:"id"`
	Comment         string                `json:"comment"`
	Detail          string                `json:"detail"`
	PayType         string                `json:"paytype"`
	Require         int                   `json:"require"`
	Status          string                `json:"status"`
	CreateAt        string                `json:"create_at"`
	PassAt          string                `json:"pass_at"`
	OnAt            string                `json:"on_at"`
	ContactOperator string                `json:"contact_operator"`
	ContactEngineer string                `json:"contact_engineer"`
	LocationInfo    location.LocationInfo `json:"location_info"`
	Msg             string                `json:"message"`
}

var inwhich = []string{"comment", "id"}

type ServiceNew struct {
	Comment string `json:"comment"`
	Detail  string `json:"detail"`
	PayType string `json:"paytype"`
	LocId   uint   `json:"location"`
	Require int    `json:"require"`
	ReCAPTCHA	string `json:"g_recaptcha_response"`
}

func Service(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodPost:
		NewService(c)
	case http.MethodGet:
		GetServiceByFilter(c)
	default:
		resp.ERROR(c, http.StatusMethodNotAllowed, constant.MsgMethodNotAllowed)
		return
	}
}

// @Tags		Service
// @Summary		Create New service
// @Description	with jwt auth, new a service with informations
// @Description	**must fill:**
// @Description <pre>{
// @Description   comment, // for user to remember, 0~10 chinese
// @Description   detail, // for operator to check
// @Description   paytype, // = enum('month','year')
// @Description   location, // = location id
// @Description   require, // 0~7, 100 for private, 010 for client, 001 for test
// @Description }</pre>
// @Description	*can be null:* {
// @Description }
// @Accept		json
// @Product		json
// @Param		g-recaptcha-response query string true "response token got from google reCAPTCHA (expected action: create_service)"
// @Param		Authorization header string true "Bearer [access_token]"
// @Param		body body ServiceNew true "service information"
// @Success		200 {object} resp.Response "ok"
// @Failure		400 {object} resp.Response "reCAPTCHA failed<br>format error<br>field empty<br>permission denied<br>invalid info"
// @Router		/service [POST]
func NewService(c *gin.Context) {
	var service ServiceNew
	if err := c.ShouldBindJSON(&service); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if reCAPTCHASucceed, err := utils.ReCAPTCHA(service.ReCAPTCHA, "create_service", c.ClientIP()); !reCAPTCHASucceed {
		log.Println(err)
		resp.ERROR(c, http.StatusForbidden, "reCAPTCHA failed")
		return
	}
	if service.Comment == "" ||
		service.Detail == "" ||
		service.PayType == "" {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFieldEmpty)
		return
	}
	if !model.ValidCommentLen(service.Comment) ||
		!model.ValidDetailLen(service.Detail) ||
		!model.ValidPayType(service.PayType) ||
		!model.ValidRequire(service.Require) {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	var loc model.Location
	if err := model.DB.First(&loc, service.LocId).Error; err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}

	var u model.User
	model.DB.First(&u, c.MustGet("userId").(uint))

	if u.Role != constant.RoleCustomer {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgPermissionDenied)
		return
	}

	// assign the new service to a random engineer and a operator
	engineer := getRandomUser(constant.RoleEngineer)
	operator := getRandomUser(constant.RoleOperator)

	s := model.Service{
		Comment:    service.Comment,
		Detail:     service.Detail,
		PayType:    service.PayType,
		Status:     constant.StatusWaiting,
		LocationID: service.LocId,
		Require:    service.Require,
		Users: []*model.User{
			&u, &engineer, &operator,
		},
	}
	if err := model.DB.Create(&s).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, err.Error())
	} else {
		resp.OK(c, http.StatusOK, gin.H{})
	}
}

// @Tags		Service
// @Summary		Get service by filter
// @Description	with jwt auth and filter conditions, get user's services
// @Description **must fill:**
// @Description <pre> {
// @Description 	offset,
// @Description 	limit, // cannot be 0, should be <= 20
// @Descriotion 	// in range [offset, offset+limit)
// @Description 	status, // array of type needed [constant.StatusWaiting, constant.StatusPass, constant.StatusOn, constant.StatusPause, constant.StatusCanceled]
// @Description } </pre>
// @Description *can be null:*
// @Description <pre> {
// @Description 	search, // key words seperated by space, len<=40, search by OR method // support uuid & comment
// @Description } </pre>
// @Description	return :
// @Description	<pre> {
// @Description   "count": count // total number of services satisfied the condition
// @Description   "services": [{
// @Description   	id, // uuid
// @Description   	comment, // comment
// @Description   	paytype // = enum('month','year')
// @Description   	require, // 0~7, 100 for private, 010 for client, 001 for test
// @Description   	status, // = enum('waiting', 'pass', 'on', 'pause', 'canceled'); waiting： 等待运营师审核; pass： 审核完了等待工程师部署; on：正常工作中; pause：被暂停了，比如欠费; canceled：被删了
// @Description   	create_at, // in format: "YYYY-MM-DD HH:mm:ss"
// @Description		message, // extra message
// @Description   }, ] // array of services by create_at descend
// @Description	}</pre>
// @Accept		json
// @Product		json
// @Param 		Authorization header string true "Bearer [access_token]"
// @Param 		offset query int true "offset"
// @Param 		limit query int true "0<limit<=20"
// @Param 		status query []string true "list of status" collectionFormat(multi)
// @Param 		search query string false "0<limit<=40"
// @Success		200 {object} resp.Response{data=[]ServiceArray} "ok"
// @Failure		400 {object} resp.Response "format error<br>invalid info"
// @Router		/service [GET]
func GetServiceByFilter(c *gin.Context) {
	var filter ServiceFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgFormatError)
		return
	}
	if filter.Limit == 0 || filter.Limit > 20 {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	var user model.User
	var services []model.Service
	var err error
	whereClause := model.DB
	if whereClause, err = utils.GenerateStatusClause(whereClause, filter.Status); err != nil {
		resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
		return
	}
	if filter.Search != "" {
		if len(filter.Search) <= 120 {
			if whereClause, err = utils.GenerateSearchCLause(whereClause, filter.Search, inwhich); err != nil {
				resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
				return
			}
		} else {
			resp.ERROR(c, http.StatusBadRequest, constant.MsgInvalidInfo)
			return
		}
	}

	if err := model.DB.First(&user, c.MustGet("userId").(uint)).Error; err != nil {
		resp.ERROR(c, http.StatusInternalServerError, constant.MsgTokenCannotFindUser)
		return
	}
	model.DB.Model(&user).Where(whereClause).Limit(filter.Limit).Offset(filter.Offset).Order("create_at desc").
		Association(constant.TableService).Find(&services)
	var ret []ServiceAbstract
	for _, s := range services {
		ret = append(ret, ServiceAbstract{
			ID:       s.ID.String(),
			Comment:  s.Comment,
			PayType:  s.PayType,
			Require:  s.Require,
			Status:   s.Status,
			CreateAt: model.Time2String(s.CreateAt),
			Msg:      s.Msg,
		})
	}
	total := int(model.DB.Model(&user).Where(whereClause).Association(constant.TableService).Count())
	resp.OK(c, http.StatusOK, ServiceArray{
		Count:    total,
		Services: ret,
	})
}

func getRandomUser(role string) (u model.User) {
	var count int64
	model.DB.Model(&model.User{}).Where("role = ?", role).Count(&count)
	model.DB.Offset(rand.Intn(int(count))).Where("role = ?", role).First(&u)
	return
}
